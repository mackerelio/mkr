use strict;
use warnings;
use utf8;

use Test::More;
use Test::Mock::Guard;
use File::Temp "tempfile";

use Releng::File ();

subtest "slurp" => sub {
    my ($fh, $filename) = tempfile;
    my $expected = "1 2 3 4 5";
    print $fh $expected;
    close $fh;
    is(Releng::File::slurp($filename), $expected);
};

subtest "slurp" => sub {
    my ($fh, $filename) = tempfile;
    close $fh;
    my $original = "1 2 3 4 5";
    my $expected = "1 2 3 4 5\n";
    Releng::File::spew($filename, $expected);
    is(Releng::File::slurp($filename), $expected);
};

subtest "replace" => sub {
    my ($fh, $filename) = tempfile;
    close $fh;
    my $original = "1 2 3 4 5";
    my $expected = "1 3 2 4 5\n";
    Releng::File::spew($filename, $original);
    Releng::File::replace $filename => sub {
        my $content = shift;
        $content =~ s/2\s3/3 2/;
        $content;
    };
    is(Releng::File::slurp($filename), $expected);
};

subtest "replace_if_exist" => sub {
    {
        my $called = 0;
        my $guard = mock_guard('Releng::File', {
            replace => sub { $called = 1 },
        });
        my ($fh, $filename) = tempfile;
        close $fh;
        Releng::File::replace_if_exist $filename => sub {};
        ok $called, "existing file";
    }

    {
        my $called = 0;
        my $guard = mock_guard('Releng::File', {
            replace => sub { $called = 1 },
        });
        my $filename = "awesomefile";
        Releng::File::replace_if_exist $filename => sub {};
        ok !$called, "not existing file";
    }
};

done_testing;
