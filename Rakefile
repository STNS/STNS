require 'erb'
require 'rake'
require 'rspec/core/rake_task'

task :default => "repo"

desc "run server spec"
task :spec    => 'spec:all'

task "clean_all" do
  sh "rm -rf binary/*"
  sh "rm -rf releases/*"
end

desc "delete binarys"
task "clean_bin" do
  sh "rm -rf binary/stns.bin"
end

desc "delete packages"
task "clean_pkg" do
  sh "find binary/* | grep -e 'rpm$' -e 'deb$' | xargs rm -rf"
end

desc "make package all architecture"
task "make_pkg" => %W(
  clean_pkg
  make_pkg_x86
  make_pkg_i386
)

desc "test package all architecture"
task "test_pkg" => %W(
  make_pkg
  test_pkg_x86
  test_pkg_i386
)


%w(x86 i386).each do |arch|
  desc "make package #{arch}"
  task "make_pkg_#{arch}" => %W(
    clean_bin
    ubuntu_build_#{arch}
    centos_pkg_#{arch}
    ubuntu_pkg_#{arch}
  )

  desc "test package #{arch}"
  task "test_pkg_#{arch}" => %W(
    centos_ci_#{arch}
    ubuntu_ci_#{arch}
  )
end

[
  {
    os: "centos",
    arch: %w(x86 i386),
    pkg_arch: %w(x86_64 i386)
  },
  {
    os: "ubuntu",
    arch: %w(x86 i386),
    pkg_arch: %w(amd64 i386)
  }
].each do |h|

  h[:arch].each_with_index do |arch,index|
    task "#{h[:os]}_build_#{arch}" do
      docker_run(h[:os], arch, "build")
    end unless h[:os] == "centos"

    task "#{h[:os]}_pkg_#{arch}" do
      docker_run(h[:os], arch, "pkg", h[:pkg_arch][index])
    end

    task "#{h[:os]}_ci_#{arch}" do
      docker_run(h[:os], arch, "ci", h[:pkg_arch][index])
    end
  end
end

task "make_client" do
  sh "cd ../libnss_stns && bundle exec rake clean_all && bundle exec rake make_pkg"
end

desc "make repositry"
task "repo" => %i(
  clean_all
  make_client
  make_pkg_x86
  make_pkg_i386
) do
  sh "cp -pr ../libnss_stns/binary/*.rpm binary"
  sh "cp -pr ../libnss_stns/binary/*.deb binary"

  raise 'package not found' unless %w(stns libnss-stns).all? do |f|
    sh "test -e binary/#{f}*x86_64.rpm"
    sh "test -e binary/#{f}*amd64.deb"
    sh "test -e binary/#{f}*i386.rpm"
    sh "test -e binary/#{f}*i386.deb"
  end

  raise "can't create repo" unless %w(centos ubuntu).all? {|os| docker_run(os, nil, "repo", nil, "releases") }
end

def docker_run(os, arch, task, pkg_arch=nil, dir="binary")
  content = ERB.new(open("docker/#{os}-#{task}.erb").read).result(binding)
  open("docker/tmp/#{os}-#{arch}-#{task}","w") {
    |f| f.write(content)
  }

  sh "docker build --rm -f docker/tmp/#{os}-#{arch}-#{task} -t stns:stns ."
  sh "docker run --rm -e ARCH=#{pkg_arch} --rm -it -v \"$(pwd)\"/#{dir}:/go/src/github.com/STNS/STNS/#{dir} -t stns:stns"
end

namespace :spec do
  targets = []
  Dir.glob('./spec/*').each do |dir|
    next unless File.directory?(dir)
    target = File.basename(dir)
    target = "_#{target}" if target == "default"
    targets << target
  end

  task :all     => targets
  task :default => :all

  targets.each do |target|
    original_target = target == "_default" ? target[1..-1] : target
    desc "Run serverspec tests to #{original_target}"
    RSpec::Core::RakeTask.new(target.to_sym) do |t|
      ENV['TARGET_HOST'] = original_target
      t.pattern = "spec/#{original_target}/*_spec.rb"
    end
  end
end
