require 'formula'

class Pup < Formula
  homepage 'https://github.com/EricChiang/pup'
  version '0.3.3'

  if Hardware.is_64_bit?
    url 'https://github.com/EricChiang/pup/releases/download/v0.3.3/pup_darwin_amd64.zip'
    sha1 'e5a74c032abd8bc81e4a12b06d0c071343811949'
  else
    url 'https://github.com/EricChiang/pup/releases/download/v0.3.3/pup_darwin_386.zip'
    sha1 'cd7d18cae7d8bf6af8bdb04c963156a1b217dfcb'
  end

  def install
    bin.install 'pup'
  end
end
