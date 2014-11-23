require 'formula'

class Pup < Formula
  homepage 'https://github.com/EricChiang/pup'
  version '0.3.5'

  if Hardware.is_64_bit?
    url 'https://github.com/EricChiang/pup/releases/download/v0.3.5/pup_darwin_amd64.zip'
    sha1 '6991dc9408e02adfa0ed5866eb7e284a94d79a77'
  else
    url 'https://github.com/EricChiang/pup/releases/download/v0.3.5/pup_darwin_386.zip'
    sha1 'ec58d15a39ab821caa5f903035862690bbeb4dfe'
  end

  def install
    bin.install 'pup'
  end
end
