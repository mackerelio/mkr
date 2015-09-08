package Releng::File;

use 5.014;
use strict;
use warnings;
use utf8;
use Carp;

use Exporter "import";
our @EXPORT = qw( slurp spew replace replace_if_exist );

sub slurp {
    my $file = shift;
    local $/;
    open my $fh, '<:encoding(UTF-8)', $file or die $!;
    <$fh>
}

sub spew {
    my ($file, $data) = @_;
    open my $fh, '>:encoding(UTF-8)', $file or die $!;
    $data .= "\n" if $data !~ /\n\z/ms;
    print $fh $data;
}

sub replace {
    my ($file, $code) = @_;
    my $content = $code->(slurp($file));
    spew($file, $content);
}

sub replace_if_exist {
    my ($file, $code) = @_;
    replace($file => $code) if -e $file;
}

1;
