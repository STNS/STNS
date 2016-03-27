task :default => "repo"

task "clean_all" do
  sh "rm -rf binary/*"
  sh "rm -rf releases/*"
end

task "clean_bin" do
  sh "ls -d binary/* | grep -v -e 'rpm$' -e 'deb$' | xargs rm -rf"
end

[
  %w(x86 x86_64 amd64),
  %w(i386 i386 i386)
].each do |r|
  task "build_#{r[0]}" => [:clean_bin]  do
    docker_run "ubuntu-#{r[0]}-build"
  end

  task "pkg_#{r[0]}" => ["build_#{r[0]}".to_sym] do
    sh "ls -d binary/* | grep -e '#{r[1]}.rpm$' -e '#{r[2]}.deb$'| xargs rm -rf"
    docker_run("centos-#{r[0]}-rpm", r[1])
    docker_run("ubuntu-#{r[0]}-deb", r[2])

    # check package
    sh "test -e binary/*#{r[1]}.rpm"
    sh "test -e binary/*#{r[2]}.deb"
  end
end

task "make_client" do
  sh "cd ../lib-stns && bundle exec rake make"
end

task "repo" => [:clean_all, :make_client, :pkg_i386, :pkg_x86] do
  sh "cp -pr ../lib-stns/binary/*.rpm binary"
  sh "cp -pr ../lib-stns/binary/*.deb binary"

  raise 'package not found' unless %w(stns lib-stns).all? do |f|
    sh "test -e binary/#{f}*x86_64.rpm"
    sh "test -e binary/#{f}*amd64.deb"
    sh "test -e binary/#{f}*i386.rpm"
    sh "test -e binary/#{f}*i386.deb"
  end

  raise "can't create repo" unless %w(centos ubuntu).all? {|o| docker_run("#{o}-x86-repo", "", "releases") }
end

def docker_run(file, arch="x86_64", dir="binary")
  sh "docker build --no-cache --rm -f docker/#{file} -t stns:stns ."
  sh "docker run  -e ARCH=#{arch} -it -v \"$(pwd)\"/#{dir}:/go/src/github.com/STNS/STNS/#{dir} -t stns:stns"
end
