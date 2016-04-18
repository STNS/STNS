require 'serverspec'

set :backend, :exec

def i386?
  Specinfra.backend.run_command("gcc -v 3>&2 2>&1 1>&3 | grep -e '--build=i386' -e '--build=i686'").exit_status == 0
end
