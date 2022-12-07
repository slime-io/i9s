#!/bin/bash -e
set -e

export KUBECONFIG=${KUBECONFIG:-"$HOME/.kube/config"}

# use can specify tag， like 'tag=v0.0.7-i9s ./install.sh'

run() {
  tag=${tag:-$(curl https://api.github.com/repos/slime-io/i9s/releases/latest -s|grep tag_name|sed 's/.*tag_name": "//g; s/",.*//g')}
  echo "docker run -it -d --net=host -v $KUBECONFIG:/root/.kube/config slimeio/i9s:$tag"
  containerID=$(docker run -it -d --net=host -v $KUBECONFIG:/root/.kube/config slimeio/i9s:$tag)
  echo "$containerID"
  docker cp $containerID:/bin/i9s /tmp/i9s/
  docker cp $containerID:/bin/istioctl /tmp/i9s/
  docker cp $containerID:/usr/local/bin/fx /tmp/i9s/
  docker kill $containerID
}

has() {
  type "$1" > /dev/null 2>&1
}

pre_check() {
  # mkdir dir
  if [ ! -d "/tmp/i9s" ]; then
    echo "mkdir /tmp/i9s"
    mkdir /tmp/i9s
  fi

  # check kubectl
  if ! has "kubectl"; then
    echo "you must install kubectl by yourself"
    exit 1
  fi

  # check os
  local OS
  if [[ `cat /etc/os-release 2>/dev/null` =~ CentOS ]]; then
    OS="CentOS"
  elif [[ `cat /etc/os-release 2>/dev/null` =~ Ubuntu|Debian ]]; then
    OS="DEBIAN"
  fi

  # check jq
  if [ "$OS" == "DEBIAN" ]; then
    if ! has "jq" ; then
      echo "jq is not found, prepare to install"
      apt install -y jq
    fi
    if ! has "less" ; then
      echo "less is not found, prepare to install"
      apt install -y less
    fi
  elif [ "$OS" == "CentOS"  ]; then
    if ! has "jq" ; then
      echo "jq is not found, prepare to install"
      yum install -y jq
    fi
    if ! has "less" ; then
      echo "less is not found, prepare to install"
      yum install -y less
    fi
  fi
}

check() {
  echo "download i9s..."
  run
  mv /tmp/i9s/i9s /usr/bin/
  if ! has "istioctl"; then
  echo "download istioctl..."
  mv /tmp/i9s/istioctl /usr/bin
  fi
  if ! has "fx"; then
  echo "download fx..."
  mv /tmp/i9s/fx /usr/bin
  fi
}

download_istioctl() {
  if [ ! -f "/tmp/i9s/istioctl-1.12.0-alpha.0-linux-amd64.tar.gz" ]; then
    echo "download istioctl..."
    wget --no-check-certificate -nc --directory-prefix "/tmp/i9s" "https://github.com/istio/istio/releases/download/1.12.0-alpha.0/istioctl-1.12.0-alpha.0-linux-amd64.tar.gz"
    tar zxvf /tmp/i9s/istioctl-1.12.0-alpha.0-linux-amd64.tar.gz -C /tmp/i9s/
    echo "mv istioctl to PATH..."
    mv /tmp/i9s/istioctl /usr/bin/
  fi
}

install_i9s() {
  if ! has "i9s"; then
    echo "download i9s..."
    # TODO wget --no-check-certificate -nc --directory-prefix "/tmp/i9s" "xx"
  fi
}

exec 2>&1
echo "i9s install..."
pre_check
check
echo "succeed"
