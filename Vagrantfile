# -*- mode: ruby -*-
# vi: set ft=ruby :
#

LINUX_BASE_BOX = "bento/ubuntu-19.10"

LINUX_IP_ADDRESS = "10.199.0.200"

Vagrant.configure(2) do |config|
  config.vm.define "eoan", autostart: true, primary: true do |vmCfg|
	vmCfg.vm.box = LINUX_BASE_BOX
	vmCfg.vm.hostname = "linux"

    vmCfg.vm.provider "virtualbox" do |v|
	  v.customize ["modifyvm", :id, "--cableconnected1", "on", "--audio", "none"]
	  v.customize ["modifyvm", :id, "--cableconnected1", "on"]
	  v.memory = 4096
	  v.cpus = 2
    end

    vmCfg.vm.provision "shell",
		               privileged: true,
		               inline: 'rm -f /home/vagrant/linux.iso'

    vmCfg.vm.provision "shell",
		               privileged: true,
		               path: './provisioning/vagrant-linux-priv-go.sh'

    vmCfg.vm.provision "shell",
		               privileged: true,
		               path: './provisioning/vagrant-linux-priv-config.sh'

    vmCfg.vm.provision "shell",
		               privileged: true,
		               path: './provisioning/vagrant-linux-priv-docker.sh'

    vmCfg.vm.provision "shell",
		               privileged: true,
		               path: './provisioning/vagrant-linux-priv-consul.sh'

    vmCfg.vm.provision "shell",
		               privileged: true,
		               path: './provisioning/vagrant-linux-priv-bpf.sh'

    vmCfg.vm.provision "shell",
		               privileged: false,
		               path: './provisioning/vagrant-linux-unpriv-bpf.sh'

    vmCfg.vm.provision "shell",
		               privileged: false,
		               path: './provisioning/vagrant-linux-unpriv-profile.sh'

	vmCfg.vm.synced_folder '.',
			               '/opt/gopath/src/golang-bpf'

    # Nomad and Consul
	vmCfg.vm.network :forwarded_port, guest: 4646, host: 4646, auto_correct: true
	vmCfg.vm.network :forwarded_port, guest: 8500, host: 8500, auto_correct: true
	vmCfg.vm.network :forwarded_port, guest: 4201, host: 4201, auto_correct: true
	vmCfg.vm.network :forwarded_port, guest: 49153, host: 49153, auto_correct: true

    # misc web stuff
	vmCfg.vm.network :forwarded_port, guest: 8080, host: 8080, auto_correct: true
  end
end
