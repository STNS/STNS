require 'erb'
task :default => "repo"

task "clean_all" do
  sh "rm -rf binary/*"
  sh "rm -rf releases/*"
end

task "clean_bin" do
  sh "find binary/* | grep -v -e 'rpm$' -e 'deb$' | xargs rm -rf"
end

[
  %w(x86 x86_64 amd64),
  %w(i386 i386 i386)
].each do |r|
  arch = r[0]
  arch_rpm = r[1]
  arch_deb = r[2]

  task "build_#{arch}" => [:clean_bin]  do
    content = ERB.new(open("docker/ubuntu-build.erb").read).result(binding)
    open("docker/tmp/ubuntu-#{arch}-build","w") {
      |f| f.write(content)
    }
    docker_run "tmp/ubuntu-#{arch}-build"
  end

  task "pkg_#{arch}" => ["build_#{arch}".to_sym] do
    [
      ["centos", arch_rpm, "rpm"],
      ["ubuntu", arch_deb, "deb"]
    ].each do |o|
      content = ERB.new(open("docker/#{o[0]}-pkg.erb").read).result(binding)
      open("docker/tmp/#{o[0]}-#{arch}-pkg","w") {
        |f| f.write(content)
      }

      sh "find binary/* | grep -e '#{o[1]}.#{o[2]}$' | xargs rm -rf"

      docker_run("tmp/#{o[0]}-#{arch}-pkg", o[1])
      # check package
      sh "test -e binary/*#{o[1]}.#{o[2]}"
    end
  end
end

task "make_client" do
  sh "cd ../libnss_stns && bundle exec rake make"
end

task "repo" => [:clean_all, :make_client, :pkg_i386, :pkg_x86] do
  sh "cp -pr ../libnss_stns/binary/*.rpm binary"
  sh "cp -pr ../libnss_stns/binary/*.deb binary"

  raise 'package not found' unless %w(stns libnss-stns).all? do |f|
    sh "test -e binary/#{f}*x86_64.rpm"
    sh "test -e binary/#{f}*amd64.deb"
    sh "test -e binary/#{f}*i386.rpm"
    sh "test -e binary/#{f}*i386.deb"
  end

  raise "can't create repo" unless %w(centos ubuntu).all? {|o| docker_run("#{o}-repo", "", "releases") }
end

def docker_run(file, arch="x86_64", dir="binary")
  sh "docker build --no-cache --rm -f docker/#{file} -t stns:stns ."
  sh "docker run  -e ARCH=#{arch} -it -v \"$(pwd)\"/#{dir}:/go/src/github.com/STNS/STNS/#{dir} -t stns:stns"
end
