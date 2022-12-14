module melato.org/jsoncall/example

go 1.19

replace (
	melato.org/command => ../../command
	melato.org/jsoncall => ../
)

require (
	melato.org/command v0.0.0-00010101000000-000000000000
	melato.org/jsoncall v0.0.0-00010101000000-000000000000
)

require gopkg.in/yaml.v2 v2.4.0 // indirect
