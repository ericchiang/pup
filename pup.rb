require 'formula'

class Pup < Formula
  homepage 'https://github.com/EricChiang/pup'
  version '0.3.2'

  if Hardware.is_64_bit?
    url 'https://github.com/EricChiang/pup/releases/download/v0.3.2/pup_darwin_amd64.zip'
    sha1 '9d5ad4c0b78701b1868094bf630adbbd26ae1698'
  else
    url 'https://github.com/EricChiang/pup/releases/download/v0.3.2/pup_darwin_386.zip'
    sha1 '21487bc5abdac34021f25444ab481e267bccbd72'
  end

  def install
    bin.install 'pup'
  end
end
