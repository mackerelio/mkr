package Releng::Utils;

use 5.014;
use strict;
use warnings;
use utf8;
use Carp;

use Exporter "import";
our @EXPORT = qw( scope_guard );

package __g {
    sub new {
        my ($class, $code) = @_;
        bless $code, $class;
    }
    sub DESTROY {
        my $self = shift;
        $self->();
    }
}

sub scope_guard(&) {
    my $code = shift;
    __g->new($code);
}

1;
