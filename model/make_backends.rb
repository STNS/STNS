require 'erb'

getter_user_groups = [
	["FindUserByID", "int"],
	["FindUserByName", "string"],
	["FindGroupByID", "int"],
	["FindGroupByName", "string"],
	["Users", nil],
	["Groups", nil]
]

getter_user_groups_template = <<~EOS
func (gb Backends) <%= t[0] %>(<%= t[1] ? "v " + t[1] : "" %>) (map[string]UserGroup, error) {
	r := map[string]UserGroup{}
  var notfound error
  eg := errgroup.Group{}
	for _, b := range gb {
    eg.Go(func() error {
      lr, err := b.<%= t[0] %>(<%= t[1] ? "v" : "" %>)
      if err != nil {
        switch err.(type) {
          case NotFoundError:
            notfound = err
          default:
            return err
        }
      }
      r = mergeUserGroup(r, lr)
      return nil
    })
	}
  if err := eg.Wait(); err != nil {
      return nil, err
  }

  // record notfound
  if len(r) == 0 {
    return nil, notfound
  }

	return r, nil
}
EOS

highlow = %w(
  HighestUserID
  LowestUserID
  HighestGroupID
  LowestGroupID
)

highlow_template = <<~EOS
func (gb Backends) <%= t %>() int {
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


setter_user_groups = [
	["CreateUser", "v UserGroup"],
	["CreateGroup", "v UserGroup"],
	["DeleteUser", "v int"],
	["DeleteGroup", "v int"],
	["UpdateUser", "v int", "vv UserGroup"],
	["UpdateGroup", "v int", "vv UserGroup"],
]

setter_user_groups_template = <<~EOS
func (gb Backends) <%= t[0] %>(<%= t[1] %><%= t[2] ? "," + t[2] : "" %>) error {
  eg := errgroup.Group{}
	for _, b := range gb {
    eg.Go(func() error {
      err := b.<%= t[0] %>(<%= t[1] ? "v" : "" %><%= t[2] ? ", vv" : "" %>)
      if err != nil {
          return err
      }
      return nil
    })
	}
  if err := eg.Wait(); err != nil {
      return nil
  }

	return nil
}
EOS

fname = 'model/backends.go'
file = File.open(fname,'w')

file.puts "package model"

getter_user_groups.each do |t|
  erb = ERB.new(getter_user_groups_template)
  file.puts erb.result(binding)
end

setter_user_groups.each do |t|
  erb = ERB.new(setter_user_groups_template)
  file.puts erb.result(binding)
end

highlow .each do |t|
  erb = ERB.new(highlow_template)
  file.puts erb.result(binding)
end
file.close

`go fmt #{fname}`
`goimports -w #{fname}`
