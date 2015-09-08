use strict;
use warnings;
use utf8;

use Test::More;

use Releng::Utils ();

subtest "guarded" => sub {
    my $old_value = 1;
    my $new_value = 5;
    my $value = $old_value;
    {
        $value = $new_value;
        my $g = Releng::Utils::scope_guard {
            $value = $old_value;
        };
        is($value, $new_value, "inner scope");
    };
    is($value, $old_value, "outer scope");
};

done_testing;
