%global gopath %{_tmppath}/gopath
%global gopup %{gopath}/src/github.com/ericchiang/pup

Name:       pup
Version:    0.4.0
Release:    1%{?dist}
Summary:    CLI for processing HTML

License:    MIT
URL:        https://github.com/ericchiang/pup
Source0:    https://github.com/ericchiang/pup/archive/v%{version}.tar.gz

BuildRequires:  golang

%description
Command-line tool for processing HTML

%prep
%setup -q -n %{name}-%{version}
rm -rf %{gopath}
mkdir -p %{gopath} %{gopath}/pkg %{gopath}/bin %{gopath}/src %{gopup}

%build
cp -R * %{gopup}
cd %{gopup}
GOPATH=%{gopath} go get -v .
#cp -a `pwd`/vendor/* %{gopath}/src
GOPATH=%{gopath} go build -a -ldflags '-s -w -extldflags "-static"' .

%install
install -d %{buildroot}%{_bindir}
install -m 755 %{gopup}/%{name} %{buildroot}%{_bindir}

%files
%doc LICENSE
%doc README.md
%{_bindir}/%{name}


%changelog
* Tue Jul 14 2020 Kilian Cavalotti <kilian@stanford.edu> 0.4.0-1
- initial RPM release
