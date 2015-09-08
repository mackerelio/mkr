package Releng::Logger;

use 5.014;
use strict;
use warnings;
use utf8;

use Exporter "import";
our @EXPORT = qw( infof warnf debugf errorf );

use Term::ANSIColor qw(colored);
use constant { LOG_DEBUG => 1, LOG_INFO => 2, LOG_WARN => 3, LOG_ERROR => 4 };

sub DEBUG() { $ENV{MC_RELENG_DEBUG} }

my $Colors = {
    LOG_DEBUG, => "green",
    LOG_WARN,  => "yellow",
    LOG_INFO,  => "cyan",
    LOG_ERROR, => "red",
};

sub _printf {
    my $type = pop;
    return if $type == LOG_DEBUG && !DEBUG;
    my ($temp, @args) = @_;
    my $msg = sprintf($temp, map { defined($_) ? $_ : "-" } @args);
    $msg = colored $msg, $Colors->{$type} if defined $type;
    my $fh = $type && $type >= LOG_WARN ? *STDERR : *STDOUT;
    print $fh $msg;
}

sub infof  {_printf(@_, LOG_INFO)}
sub warnf  {_printf(@_, LOG_WARN)}
sub debugf {_printf(@_, LOG_DEBUG)}
sub errorf {
    my(@msg) = @_;
    _printf(@msg, LOG_ERROR);

    my $fmt = shift @msg;
    die sprintf($fmt, @msg);
}

1;
