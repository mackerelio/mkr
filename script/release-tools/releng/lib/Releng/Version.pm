package Releng::Version;

use 5.014;
use strict;
use warnings;
use utf8;
use version;

use Exporter "import";
our @EXPORT = qw( suggest_next_version is_valid_version last_version );

sub suggest_next_version {
    my ($major, $minor, $patch) = @{shift->{version}};
    version->parse(join '.', $major, ++$minor, 0);
}

sub is_valid_version {
    my $version = shift;
    ($version =~ /^\d+\.\d+\.\d+$/) or return 0;
    eval { version->parse($version) };
    ($@ ? 0 : 1);
}

sub last_version {
    my @tags = @_;

    my ($tag) =
        sort { version->parse($b) <=> version->parse($a) }
        map {/^v([0-9]+(?:\.[0-9]+)+)$/; $1 || ()}
        map {chomp; $_} @tags;
    version->parse($tag);
}

1;
