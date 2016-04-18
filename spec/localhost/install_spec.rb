require 'spec_helper'

describe file("/etc/stns/stns.conf") do
  it { should be_mode(644) }
  it { should be_owned_by('root') }
  it { should be_grouped_into('root') }
end

describe file("/usr/local/bin/stns") do
  it { should be_owned_by('root') }
  it { should be_grouped_into('root') }
end

describe command("file /usr/local/bin/stns") do
  bit = i386? ? "32" : "64"
  its(:stdout) { should match /#{bit}-bit/ }
end

%w(
  start
  restart
  reload
  checkconf
).each do |cmd|
  describe command("service stns #{cmd}") do
    its(:exit_status) { should eq 0 }
  end
end

describe service('stns') do
  it { should be_enabled }
  it { should be_running }
end
