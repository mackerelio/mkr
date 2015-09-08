package Releng;

use 5.014;
use strict;
use warnings;
use utf8;
use Carp;

use HTTP::Tiny;
use JSON::PP;
use ExtUtils::MakeMaker qw/prompt/;
use Time::Piece;
use POSIX qw(setlocale LC_TIME);
use version;

use Releng::File;
use Releng::Utils;
use Releng::Version;
use Releng::Logger;

sub new {
    my $class = shift;

    my $repo_name = $class->_extract_repo_name;
    bless {repo_name => $repo_name}, $class;
}

sub _extract_repo_name {
    chomp(my $url = `hub browse -u`);
    $url =~ /github.com/ or croak("not github project");
    $url =~ /mackerelio/ or croak("not mackerelio project");
    [split "/", $url]->[-1];
}

sub repo_name {
    shift->{repo_name};
}

sub last_release {
    last_version `git tag`;
}

sub decide_next_version {
    my $self = shift;

    my $current_version = shift;
    my $suggest_version = suggest_next_version($current_version);

    my $next_version_str = prompt("next version", $suggest_version);
    say $next_version_str;

    if (!is_valid_version($next_version_str)) {
        croak(qq{"$next_version_str" is invalid version string\n});
    }

    my $next_version = version->parse($next_version_str);

    if ($next_version < $current_version) {
        croak(qq{"$next_version_str" is smaller than current version "$next_version_str"\n});
    }

    $next_version;
}

sub merged_prs {
    my $self = shift;

    my $current_tag = shift;
    my @pull_nums = sort {$a <=> $b} map {m/Merge pull request #([0-9]+) /; $1 || ()  } `git log v$current_tag... --merges --oneline`;

    my @releases;
    my $ua = HTTP::Tiny->new;
    for my $pull_num (@pull_nums) {
        my $url = sprintf "https://api.github.com/repos/mackerelio/%s/pulls/%d?state=closed", $self->{repo_name}, $pull_num;
        my $res = $ua->get($url);
        unless ($res->{success}) {
            warnf "request to $url failed\n";
            exit;
        }
        my $data = eval { decode_json $res->{content} };
        if ($@) {
            warnf "parse json failed. url: $url\n";
            next;
        }

        push @releases, {
            num   => $pull_num,
            title => $data->{title},
            user  => $data->{user}{login},
            url   => $data->{html_url},
        } if $data->{title} !~ /\[nit\]/i;
    }
    @releases;
}

sub build_pull_request_body {
    my $self = shift;

    my ($next_version, @releases) = @_;
    my $body = "Release version $next_version\n\n";
    for my $rel (@releases) {
        $body .= sprintf "- %s #%s\n", $rel->{title}, $rel->{num};
    }
    $body;
}

sub update_versions {
    my $self = shift;

    my ($current_version, $next_version) = @_;

    my $cur_ver_reg = quotemeta $current_version;

    my $travis_yml_path = ".travis.yml";
    replace_if_exist $travis_yml_path => sub {
        my $content = shift;
        $content =~ s/$cur_ver_reg/$next_version/msg;
        $content;
    };

    my $appveyor_yml_path = "appveyor.yml";
    replace_if_exist $appveyor_yml_path => sub {
        my $content = shift;

        my $crnt_build_ver = join ".", splice(@{$current_version->{version}}, 0, 2), "{build}";
        my $next_build_ver = join ".", splice(@{$next_version->{version}},    0, 2), "{build}";

        my $crnt_build_ver_reg = quotemeta($crnt_build_ver);

        $content =~ s/(version:\s+)$crnt_build_ver_reg/$1$next_build_ver/msg;
        $content;
    };

    my $repo_name = $self->repo_name;
    my $rpm_spec_file_path = "packaging/rpm/$repo_name.spec";
    replace_if_exist $rpm_spec_file_path => sub {
        my $content = shift;
        $content =~ s/^(Version:\s+)$cur_ver_reg/$1$next_version/ms;
        $content;
    };
}

sub update_changelog {
    my $self = shift;

    my ($next_version, @releases) = @_;

    my $repo_name = $self->repo_name;
    chomp(my $email = `git config user.email`);
    chomp(my $name  = `git config user.name`);

    my $old_locale = setlocale(LC_TIME);
    setlocale(LC_TIME, "C");
    my $g = scope_guard {
        setlocale(LC_TIME, $old_locale);
    };

    my $now = localtime;

    my $deb_changelog_path = "packaging/deb/debian/changelog";
    replace_if_exist $deb_changelog_path => sub {
        my $content = shift;

        my $update = "$repo_name ($next_version-1) stable; urgency=low\n\n";
        for my $rel (@releases) {
            $update .= sprintf "  * %s (by %s)\n    <%s>\n", $rel->{title}, $rel->{user}, $rel->{url};
        }
        $update .= sprintf "\n -- %s <%s>  %s\n\n", $name, $email, $now->strftime("%a, %d %b %Y %H:%M:%S %z");
        $update . $content;
    };

    my $rpm_spec_file_path = "packaging/rpm/$repo_name.spec";
    replace_if_exist $rpm_spec_file_path => sub {
        my $content = shift;

        my $update = sprintf "* %s <%s> - %s-1\n", $now->strftime('%a %b %d %Y'), $email, $next_version;
        for my $rel (@releases) {
            $update .= sprintf "- %s (by %s)\n", $rel->{title}, $rel->{user};
        }
        $content =~ s/%changelog/%changelog\n$update/;
        $content;
    };

    my $change_log_path = "CHANGELOG.md";
    replace_if_exist $change_log_path => sub {
        my $content = shift;

        my $update = sprintf "\n\n## %s (%s)\n\n", $next_version, $now->strftime('%Y-%m-%d');
        for my $rel (@releases) {
            $update .= sprintf "* %s #%d (%s)\n", $rel->{title}, $rel->{num}, $rel->{user};
        }
        $content =~ s/\A# Changelog/# Changelog$update/;
        $content;
    };
}

1;
