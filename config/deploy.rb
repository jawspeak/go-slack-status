# config valid only for current version of Capistrano
lock '3.4.1'

set :application, 'slack-status-tool'
set :repo_url, 'git@example.com:me/my_repo.git'
# Default branch is :master
ask :branch, `git rev-parse --abbrev-ref HEAD`.chomp

# Default deploy_to directory is /var/www/my_app_name
set :deploy_to, '/home/jaw/slack-status-tool'

# Default value for :scm is :git
# set :scm, :git

# Default value for :format is :pretty
# set :format, :pretty

# Default value for :log_level is :debug
# set :log_level, :debug

# Default value for :pty is false
# set :pty, true

# Default value for :linked_files is []
# set :linked_files, fetch(:linked_files, []).push('config/database.yml', 'config/secrets.yml')

# Default value for linked_dirs is []
# set :linked_dirs, fetch(:linked_dirs, []).push('log', 'tmp/pids', 'tmp/cache', 'tmp/sockets', 'vendor/bundle', 'public/system')

# Default value for default_env is {}
# set :default_env, { path: "/opt/ruby/bin:$PATH" }

# Default value for keep_releases is 5
# set :keep_releases, 5

namespace :deploy do


  desc "Remove the remote user's crontab"
  task :remove_crontab do
    run "crontab -r; true" # ignore non-zero status when there is no crontab
  end

  desc "Generate and install crontabs"
  task :install_crontab do
    crontab = ERB.new(File.read("config/deploy/crontab.erb"), nil, '-').result(binding)

    tmp_crontab_path = "#{current_path}/crontab.tmp"
    put crontab, tmp_crontab_path
    run "crontab #{tmp_crontab_path} && rm -f #{tmp_crontab_path}"
  end

end
