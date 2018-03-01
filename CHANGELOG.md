# Changelog

## 0.27.0 (2018-03-01)

* Support empty threshold for monitors #148 (itchyny)


## 0.26.0 (2018-01-23)

* Fix copying plugin archives while installing a plugin on Windows #144 (itchyny)
* update rpm-v2 task for building Amazon Linux 2 package #143 (hayajo)


## 0.25.0 (2018-01-10)

* Add plugin document links #141 (shibayu36)
* introduce goxz and adjust deps #139 (Songmu)
* add appveyor.yml and adjust tests for windows #140 (Songmu)
* Define plugin default installation path in windows environment #138 (Songmu)


## 0.24.1 (2017-12-13)

* Rebuild to avoid panic when action of check-plugin was not specified #136 (astj)


## 0.24.0 (2017-12-12)

* Support maxCheckAttempts for host metric and service metric monitors #134 (itchyny)


## 0.23.0 (2017-11-28)

* [plugin.install] support direct URL target #130 (Songmu)
* [plugin.install] support tarball archives #131 (Songmu)
* fix hostId option flag in the command help of throw and metrics commands #129 (astj)
* Refactor mkr plugin install implementation #127 (shibayu36)


## 0.22.0 (2017-10-26)

* Release mkr plugin install #125 (shibayu36)
* Add metrics command #119 (edangelion)


## 0.21.0 (2017-10-04)

* Use new API BaseURL #116 (astj)


## 0.20.0 (2017-09-27)

* build with Go 1.9 #114 (astj)


## 0.19.0 (2017-09-20)

* Support fetch command to retrieve many hosts #112 (itchyny)
* Prefer apibase in mackerel-agent confFile #108 (astj)


## 0.18.0 (2017-09-12)

* add --customIdentifier option to mkr create #110 (astj)


## 0.17.0 (2017-08-23)

* [dashboards] Add unit to expression graph #106 (edangelion)
* [dashboards] Add title param to expression graph #104 (edangelion)


## 0.16.1 (2017-06-07)

* v2 packages #102 (Songmu)


## 0.16.0 (2017-05-09)

* Add services subcommand #97 (yuuki)


## 0.15.0 (2017-04-06)

* bump Go to 1.8 #95 (astj)


## 0.14.5 (2017-03-27)

* Colors on Windows #93 (mattn)


## 0.14.4 (2017-03-22)

* use new bot token #88 (daiksy)
* use new bot token #89 (daiksy)
* Workaround for git fetch failure #90 (daiksy)
* Apply git fetch workaround #91 (astj)


## 0.14.3 (2017-02-16)

* Support annotations command for graph annotation #83 (syou6162)
* Improve help management and fix usage help for command #85 (haya14busa)
* remove unused functions #86 (haya14busa)


## 0.14.2 (2017-02-08)

* [monitors diff] Add the "to-remote" bool flag #82 (yoheimuta)


## 0.14.1 (2017-01-11)

* formatter.NewAsciiFormatter now needs config #80 (astj)


## 0.14.0 (2016-12-21)

* Support expression monitor alerts in mkr alerts list #78 (itchyny)


## 0.13.0 (2016-11-29)

* remove unreachable code: monitor type cannot be "check" #72 (haya14busa)
* Fix the links to the api documents #73 (itchyny)
* catch up monitor interface changes of mackerel-client-go #74 (haya14busa)
* Introduce yudai/gojsondiff for `mkr monitors diff` #75 (haya14busa)
* fix test according to mackerel-client-go changes #76 (haya14busa)


## 0.12.0 (2016-10-27)

* Rename a dependent package #68 (usk81)
* Support `-apibase` option #69 (astj)
* [breaking change] Prepend `custom.` prefix to host metric name by default #70 (astj)


## 0.11.3 (2016-07-14)

* fix `validateRules()`,  when monitor has rule of "expression". #66 (daiksy)


## 0.11.2 (2016-06-23)

* replace angle brackets for json #63 (daiksy)


## 0.11.1 (2016-06-10)

* fix version number #61 (stanaka)


## 0.11.0 (2016-06-09)

* add dashboard generator #56 (daiksy)
* Add flag to overwrite host's roles  #58 (haya14busa)


## 0.10.1 (2016-05-25)

* fix signnatures. codegangsta/cli #54 (tknzk)


## 0.10.0 (2016-05-10)

* support `isMute` field of monitors #49 (Songmu)
* support boolean at isEmpty #51 (stanaka)
* bump up go version to 1.6.2 #52 (stanaka)


## 0.9.1 (2016-03-25)

* use GOARCH=amd64 for now #41 (Songmu)


## 0.9.0 (2016-02-18)

* Support displayName of host's json #39 (stanaka)


## 0.8.1 (2016-01-07)

* fix handling host-status option #37 (stanaka)


## 0.8.0 (2016-01-06)

* support alerts subcommand #31 (stanaka)
* Fix README example about mkr throw #32 (yuuki1)
* Build with Go 1.5 #33 (itchyny)
* Fixed the english used in the command descriptions #35 (stefafafan)


## 0.7.1 (2015-11-12)

* support `notificationIntervai` field in monitors (stanaka)
* [bug] fix json parameter s/hostID/hostId/g (Songmu)

## 0.7.0 (2015-10-26)

* append newline to the end of monitors.json #23 (Songmu)
* fix printMonitor #24 (Songmu)
* fix diff output between slices #25 (Songmu)

## 0.6.0 (2015-10-15)

* Fix update command bug about overwriting hostname #17 (y_uuki)
* Stop the parallel request sending temporarily #18 (y_uuki)
* Suppress to display empty fields when mkr monitors diff #20 (by stanaka)

## 0.5.0 (2015-09-14)

 * add fields for external URL monitors (stanaka)

## 0.4.1 (2015-08-28)

* Create deb/rpm package for Linux release #11 (Sixeight)


## 0.3.0 (2015-07-05)

* [feature] add --conf option to specify conf file path #4 (Sixeight)
* [fix] Fix update command as firstaid #7 (Sixeight)

## 0.2.0 (2015-06-18)

* [feature] add -f flag to hosts command to format the output #2 (motemen)
