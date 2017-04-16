#!/bin/bash

GOVER=1.8.1
ARCH=amd64

updatePackages(){
    sudo apt-get update
    sudo apt-get -y upgrade

    echo "deb http://httpredir.debian.org/debian trusty main" | sudo tee -a /etc/apt/sources.list.d/python-trusty.list
    echo "deb-src http://httpredir.debian.org/debian trusty main" | sudo tee -a /etc/apt/sources.list.d/python-trusty.list
    echo "deb http://httpredir.debian.org/debian trusty-updates main" | sudo tee -a /etc/apt/sources.list.d/python-trusty.list
    echo "deb-src http://httpredir.debian.org/debian trusty-updates main" | sudo tee -a /etc/apt/sources.list.d/python-trusty.list
    echo "deb http://security.debian.org/ trusty/updates main" | sudo tee -a /etc/apt/sources.list.d/python-trusty.list
    echo "deb-src http://security.debian.org/ trusty/updates main" | sudo tee -a /etc/apt/sources.list.d/python-trusty.list

    #sudo mv python-trusty.list /etc/apt/sources.list.d/python-trusty.list

    sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 8B48AD6246925553
    sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 7638D0442B90D010
    sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 9D6D8F6BC857C906

    sudo echo 'Package: *' >> python-trusty-pin
    sudo echo 'Pin: release o=Debian' >> python-trusty-pin
    sudo echo 'Pin-Priority: -10' >> python-trusty-pin
    sudo mv python-trusty-pin /etc/apt/preferences.d/python-trusty-pin
    sudo apt-get update
    sudo apt-get install -y git
}

installGo(){
    sudo curl -O https://storage.googleapis.com/golang/go${GOVER}.linux-${ARCH}.tar.gz
    sudo tar -C /usr/local -xzf go${GOVER}.linux-${ARCH}.tar.gz
    echo "export PATH=$PATH:/usr/local/go/bin" | sudo tee -a /etc/profile
    export PATH=$PATH:/usr/local/go/bin
    cd /vagrant
    echo "export GOPATH=/vagrant/vendors" | sudo tee -a /etc/profile
    export GOPATH=/vagrant/vendors
}

installDeps(){
    go get github.com/eclipse/paho.mqtt.golang
    go get golang.org/x/net/websocket
    go get github.com/go-ini/ini
}

updatePackages
installGo
#installDeps