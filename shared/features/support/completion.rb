# Driver for completion tests executed via a separate tmux pane in which we
# spawn an interactive shell, send keystrokes to and inspect the outcome of
# tab-completion.
#
# Prerequisites:
# - tmux
# - bash
# - zsh
# - fish
# - git

require 'fileutils'
require 'rspec/expectations'
require 'pathname'

tmpdir = Pathname.new(ENV.fetch('TMPDIR', '/tmp')) + 'hub-test'
cpldir = tmpdir + 'completion'
zsh_completion = File.expand_path('../../../etc/hub.zsh_completion', __FILE__)
bash_completion = File.expand_path('../../../etc/hub.bash_completion.sh', __FILE__)
fish_completion = File.expand_path('../../../etc/hub.fish_completion', __FILE__)

_git_prefix = nil

git_prefix = lambda {
  _git_prefix ||= begin
    git_core = Pathname.new(`git --exec-path`.chomp)
    git_core.dirname.dirname
  end
}

git_distributed_zsh_completion = lambda {
  [ git_prefix.call + 'contrib/completion/git-completion.zsh',
    git_prefix.call + 'share/git-core/contrib/completion/git-completion.zsh',
    git_prefix.call + 'share/zsh/site-functions/_git',
  ].detect {|p| p.exist? }
}

git_distributed_bash_completion = lambda {
  [ git_prefix.call + 'contrib/completion/git-completion.bash',
    git_prefix.call + 'share/git-core/contrib/completion/git-completion.bash',
    git_prefix.call + 'share/git-core/git-completion.bash',
    git_prefix.call + 'etc/bash_completion.d/git-completion.bash',
    Pathname.new('/etc/bash_completion.d/git'),
    Pathname.new('/usr/share/bash-completion/completions/git'),
    Pathname.new('/usr/share/bash-completion/git'),
  ].detect {|p| p.exist? }
}

link_completion = Proc.new { |from, name|
  raise ArgumentError, from.to_s unless File.exist?(from)
  FileUtils.mkdir_p(cpldir)
  FileUtils.ln_s(from, cpldir + name, :force => true)
}

create_file = lambda { |name, &block|
  FileUtils.mkdir_p(File.dirname(name))
  File.open(name, 'w', &block)
}

setup_tmp_home = lambda { |shell|
  FileUtils.rm_rf(tmpdir)

  case shell
  when 'zsh'
    create_file.call(tmpdir + '.zshrc') do |zshrc|
      zshrc.write <<-SH
        PS1='$ '
        for site_fn in /usr/{local/,}share/zsh/site-functions; do
          fpath=(${fpath#\$site_fn})
        done
        fpath=('#{cpldir}' $fpath)
        alias git=hub
        autoload -U compinit
        compinit -i
      SH
    end
  when 'bash'
    create_file.call(tmpdir + '.bashrc') do |bashrc|
      bashrc.write <<-SH
        PS1='$ '
        alias git=hub
        . '#{git_distributed_bash_completion.call}'
        . '#{bash_completion}'
      SH
    end
  when 'fish'
    create_file.call(tmpdir + '.config/fish/config.fish') do |fishcfg|
      fishcfg.write <<-SH
        function fish_prompt
          echo '$ '
        end
      SH
    end

    create_file.call(tmpdir + '.config/fish/functions/git.fish') do |gitfn|
      gitfn.write <<-SH
        function git --wraps hub
          hub $argv
        end
      SH
    end

    completion_dest = tmpdir + '.config/fish/completions/hub.fish'
    FileUtils.mkdir_p(File.dirname(completion_dest))
    FileUtils.ln_s(fish_completion, completion_dest)
  end
}

$tmux = nil
$installed_shells = Hash.new { |cache, shell|
  `which #{shell} 2>/dev/null`
  cache[shell] = $?.success?
}

Before('@completion') do
  unless $tmux
    $tmux = %w[tmux -L hub-test]
    system(*($tmux + %w[new-session -ds hub]))
    at_exit do
      system(*($tmux + %w[kill-server]))
    end
  end
end

After('@completion') do
  tmux_kill_pane
end

World Module.new {
  attr_reader :shell

  def set_shell(shell)
    if $installed_shells[shell]
      @shell = shell
      true
    else
      false
    end
  end

  define_method(:tmux_pane) do
    return @tmux_pane if tmux_pane?
    Dir.chdir(tmpdir) do
      @tmux_pane = `#{$tmux.join(' ')} new-window -dP -n test 'env HOME="#{tmpdir}" #{shell}'`.chomp
    end
  end

  def tmux_pane?
    defined?(@tmux_pane) && @tmux_pane
  end

  def tmux_pane_contents
    system(*($tmux + ['capture-pane', '-t', tmux_pane]))
    `#{$tmux.join(' ')} show-buffer`.rstrip
  end

  def tmux_send_keys(*keys)
    system(*($tmux + ['send-keys', '-t', tmux_pane, *keys]))
  end

  def tmux_send_tab
    @last_pane_contents = tmux_pane_contents
    tmux_send_keys('Tab')
  end

  def tmux_kill_pane
    system(*($tmux + ['kill-pane', '-t', tmux_pane])) if tmux_pane?
  end

  def tmux_wait_for_prompt
    num_waited = 0
    while tmux_pane_contents !~ /\$\Z/
      raise "timeout while waiting for shell prompt" if num_waited > 300
      sleep 0.01
      num_waited += 1
    end
  end

  def tmux_wait_for_completion
    num_waited = 0
    raise "tmux_send_tab not called first" unless defined? @last_pane_contents
    while tmux_pane_contents == @last_pane_contents
      if num_waited > 300
        if block_given? then return yield
        else
          raise "timeout while waiting for completions to expand"
        end
      end
      sleep 0.01
      num_waited += 1
    end
  end

  def tmux_output_lines
    tmux_pane_contents.split("\n").drop_while { |l| not l.start_with?('$') }.reject do |line|
      line.start_with?('$')
    end
  end

  def tmux_completion_menu
    tmux_wait_for_completion
    hash = {}
    if 'fish' == shell
      tmux_output_lines.each do |line|
        line.scan(/([^(]+)\((.+?)\)/).each do |flags, description|
          flags.strip.split(/\s+/).each do |flag|
            hash[flag] = description
          end
        end
      end
    else
      tmux_output_lines.each do |line|
        item, description = line.split(/ +-- +/, 2)
        hash[item] = description if description
      end
    end
    hash
  end

  def tmux_completion_menu_basic
    tmux_wait_for_completion
    if 'fish' == shell
      tmux_completion_menu.keys
    else
      tmux_output_lines.flat_map do |line|
        line.split(/\s+/)
      end
    end
  end
}

Given(/^my shell is (\w+)$/) do |shell|
  set_shell(shell) || pending
  setup_tmp_home.call(shell)
end

Given(/^I'm using ((?:zsh|git)-distributed) base git completions$/) do |type|
  link_completion.call(zsh_completion, '_hub')
  case type
  when 'zsh-distributed'
    raise "this combination makes no sense!" if 'bash' == shell
    expect((cpldir + '_git')).to_not be_exist
  when 'git-distributed'
    if 'zsh' == shell
      if git_zsh_completion = git_distributed_zsh_completion.call
        link_completion.call(git_zsh_completion, '_git')
        link_completion.call(git_distributed_bash_completion.call, 'git-completion.bash')
      else
        warn "warning: git-distributed zsh completion wasn't found; using zsh-distributed instead"
      end
    end
    if 'bash' == shell
      unless git_distributed_bash_completion.call
        raise "git-distributed bash completion wasn't found. Completion won't work."
      end
    end
  else
    raise ArgumentError, type
  end
end

When(/^I type "(.+?)" and press <Tab>$/) do |string|
  tmux_wait_for_prompt
  @last_command = string
  tmux_send_keys(string)
  tmux_send_tab
end

When(/^I press <Tab> again$/) do
  tmux_send_tab
end

Then(/^the completion menu should offer "([^"]+?)"( unsorted)?$/) do |items, unsorted|
  menu = tmux_completion_menu_basic
  if unsorted
    menu.sort!
    items = items.split(' ').sort.join(' ')
  end
  expect(menu.join(' ')).to eq(items)
end

Then(/^the completion menu should offer "(.+?)" with description "(.+?)"$/) do |item, description|
  menu = tmux_completion_menu
  expect(menu.keys).to include(item)
  expect(menu[item]).to eq(description)
end

Then(/^the completion menu should offer:/) do |table|
  menu = tmux_completion_menu
  expect(menu).to eq(table.rows_hash)
end

Then(/^the command should expand to "(.+?)"$/) do |cmd|
  tmux_wait_for_completion
  expect(tmux_pane_contents).to match(/^\$ #{cmd}$/)
end

Then(/^the command should not expand$/) do
  tmux_wait_for_completion { false }
  expect(tmux_pane_contents).to match(/^\$ #{@last_command}$/)
end
