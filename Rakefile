task :default => "repo"

desc "clean tmp directory"
task "clean" do
  sh "rm -rf binary/*"
  sh "rm -rf release/*"
end

desc "make binary 64bit"
task "build_64" => [:clean] do
  sh "docker build --no-cache --rm -t stns:stns ."
  sh "docker run -v \"$(pwd)\"/binary:/go/src/github.com/STNS/STNS/binary -t stns:stns"
end

desc "make package 64bit"
task "pkg_64" => [:build_64] do
  docker_run("rpm", "x86_64")
  docker_run("deb", "amd64")
end


desc "make binary 32bit"
task "build_32" => [:clean] do
  docker_run "build_32"
end

desc "make package 32bit"
task "pkg_32" => [:build_32] do
  docker_run("rpm", "i386")
  docker_run("deb_32", "i386")
end

desc "make repo data"
task "repo" => [:clean, :pkg_32, :pkg_64] do
  sh "cp -pr ../libnss_stns/binary/*.rpm binary"
  sh "cp -pr ../libnss_stns/binary/*.deb binary"
  docker_run("yum_repo", "", "releases")
  docker_run("apt_repo", "", "releases")
end

def docker_run(file, arch="x86_64", dir="binary")
  sh "docker build --no-cache --rm -f docker/#{file} -t stns:stns ."
  sh "docker run  -e TARGET=#{arch} -it -v \"$(pwd)\"/#{dir}:/go/src/github.com/STNS/STNS/#{dir} -t stns:stns"
end
