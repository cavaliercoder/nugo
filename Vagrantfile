# -*- mode: ruby -*-
# vi: set ft=ruby :

script = <<SCRIPT
#sudo apt-get update
#sudo apt-get install -y golang make

mkdir -p /home/vagrant/go
chown vagrant.vagrant /home/vagrant
cat > /home/vagrant/.profile <<EOL
export GOPATH=/home/vagrant/go
EOL

source /home/vagrant/.profile
cd /vagrant && make get-deps
SCRIPT

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/vivid64"
  config.vm.box_check_update = false

  config.vm.network "private_network", type: "dhcp"
  config.vm.network "forwarded_port", guest: 1105, host: 1105

  config.vm.provision "shell", inline: script
end
