package Releng::Command;

use 5.014;
use strict;
use warnings;
use utf8;
use Carp;

use Exporter "import";
our @EXPORT = qw( git hub which_git which_hub );

sub DEBUG() { $ENV{MC_RELENG_DEBUG} }

sub command {
  say("+ " . join " ", @_) if DEBUG; !system(@_) or croak $!
}

sub which_git {
    state $com = do {
        chomp(my $c = `which git`);
        croak "git command is required\n" unless $c;
        $c;
    };
}
sub git {
    unshift  @_, which_git; goto \&command
}

sub which_hub {
    state $com = do {
        chomp(my $c = `which hub`);
        unless ($c) {
            chomp($c = `which gh`);
        }
        croak "hub or gh command is required\n" unless $c;
        $c;
    };
}
sub hub {
    unshift @_, which_hub; goto \&command;
}
