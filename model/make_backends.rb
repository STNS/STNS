require 'erb'


user_groups = [
	["FindUserByID", "int"],
	["FindUserByName", "string"],
	["FindGroupByID", "int"],
	["FindGroupByName", "string"],
	["Users", nil],
	["Groups", nil]
]

highlow = %w(
  HighestUserID
  LowestUserID
  HighestGroupID
  LowestGroupID
)

user_groups_template = <<~EOS
func (gb GetterBackends) <%= t[0] %>(<%= t[1] ? "v " + t[1] : "" %>) (map[string]UserGroup, error) {
	r := map[string]UserGroup{}
	for _, b := range gb {
    lr, err := b.<%= t[0] %>(<%= t[1] ? "v" : "" %>)
    if err != nil {
      return nil, err
    }
		r = mergeUserGroup(r, lr)
	}
	return r, nil
}
EOS

highlow_template = <<~EOS
func (gb GetterBackends) <%= t %>() int {
	r := 0
	for _, b := range gb {
    lr := b.<%= t %>()
    if lr != 0 {
      r = lr
    }
	}
	return r
}
EOS

fname = 'model/backends.go'
file = File.open(fname,'w')

file.puts "package model"


user_groups.each do |t|
  erb = ERB.new(user_groups_template)
  file.puts erb.result(binding)
end

highlow .each do |t|
  erb = ERB.new(highlow_template)
  file.puts erb.result(binding)
end
file.close

`go fmt #{fname}`
