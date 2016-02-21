task :default => "build"

desc "make binary"
task "build" do
  sh "docker build --no-cache --rm -t stns:stns ."
  sh "docker run -v \"$(pwd)\"/binary:/go/src/github.com/STNS/STNS/binary -t stns:stns"
end

desc "make rpm package"
task "rpm" => [:build] do
  docker_run "rpm"
end

desc "make deb package"
task "deb" => [:build] do
  docker_run "deb"
end

desc "make yum repo data"
task "yum-repo" => [:rpm] do
  sh "cp -pr ../libnss_stns/binary/*.rpm binary"
  docker_run("yum_repo", "releases")
end

desc "make apt repo data"
task "apt-repo" => [:rpm] do
  sh "cp -pr ../libnss_stns/binary/*.rpm binary"
  docker_run("apt_repo", "releases")
end

def docker_run(file, dir="binary")
  sh "docker build --no-cache --rm -f docker/#{file} -t stns:stns ."
  sh "docker run -it -v \"$(pwd)\"/binary:/go/src/github.com/STNS/STNS/#{dir} -t stns:stns"
end
