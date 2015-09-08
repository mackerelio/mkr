use strict;
use warnings;
use utf8;

use Test::More;
use Test::Exception;

use Releng::Version ();
use version;

subtest "suggest_next_version (success)" => sub {
    my $versions = [
        ["0.1.0", "0.2.0"],
        ["0.1.2", "0.2.0"],
        ["0.3.2", "0.4.0"],
        ["1.0.4", "1.1.0"],
    ];
    for my $data (@$versions) {
        my $version = version->parse($data->[0]);
        my $expected = version->parse($data->[1]);
        is(Releng::Version::suggest_next_version($version), $expected);
    }
};

subtest "suggest_next_version (not version)" => sub {
    throws_ok { Releng::Version::suggest_next_version("0.20.2") } qr/Can't use string/, "string";
    throws_ok { Releng::Version::suggest_next_version(1014) } qr/Can't use string/, "number";
};

subtest "is_valid_version" => sub {
    my @valid_versions = (
        "0.1.0", "0.1.2", "0.3.2", "1.0.4",
    );
    for my $version (@valid_versions) {
        ok Releng::Version::is_valid_version($version), $version;
    }

    my @invalid_versions = (
        "0.1.0-1", "0.-1.2", "string", "1..0.4",
        "0.2", "0.2.", ".1.0", "",
    );
    for my $version (@invalid_versions) {
        ok !Releng::Version::is_valid_version($version), $version;
    }
};

subtest "last_version" => sub {
    my @tags = (
        "v0.2.2", "v0.4.0", "v0.11.1", "v0.0.5",
        "v0.3.0", "v0.11.2", "v0.0.2", "v0.10.5",
    );
    my $expected = version->parse("0.11.2");
    is(Releng::Version::last_version(@tags), $expected);
};

done_testing;
