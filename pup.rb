require 'formula'

class Pup < Formula
  homepage 'https://github.com/EricChiang/pup'
  version '0.3.4'

  if Hardware.is_64_bit?
    url 'https://github.com/EricChiang/pup/releases/download/v0.3.4/pup_darwin_amd64.zip'
    sha1 '5fec62701a49bfd5eaa4b9c980e9c06dcece78c6'
  else
    url 'https://github.com/EricChiang/pup/releases/download/v0.3.4/pup_darwin_386.zip'
    sha1 '1eb129c662d7e323c9b1e8f8ed3b8e28ce521434'
  end

  def install
    bin.install 'pup'
  end
end
