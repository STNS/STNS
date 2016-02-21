task :default => "build"

desc "make binary"
task "build" do
  sh "docker build --no-cache --rm -t stns:stns ."
  sh "docker run -v \"$(pwd)\"/binary:/go/src/github.com/STNS/STNS/binary -t stns:stns"
end

desc "make package"
task "pkg" => [:build] do
  docker_run "deb"
  docker_run "rpm"
end

desc "make repo data"
task "yum-repo" => [:pkg] do
  sh "cp -pr ../libnss_stns/binary/*.rpm binary"
  sh "cp -pr ../libnss_stns/binary/*.deb binary"
  docker_run("yum_repo", "releases")
  docker_run("apt_repo", "releases")
end

def docker_run(file, dir="binary")
  sh "docker build --no-cache --rm -f docker/#{file} -t stns:stns ."
  sh "docker run -it -v \"$(pwd)\"/binary:/go/src/github.com/STNS/STNS/#{dir} -t stns:stns"
end
