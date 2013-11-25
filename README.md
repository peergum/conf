conf
====

YAML Configuration Reader

Note
----

All configurations files used by the conf package should be located under _/usr/local/conf_ unless you change that folder in the _conf.go_ source file.

Usage
-----

Import the _Trulioo/conf_ package into your go source.

After the program initialisation, a _Conf_ variable will be created that holds your configuration. The configuration itself will be loaded from a YAML file with the name of your program, followed by "-" and the value of the -env parameter passed to it (by default "local"), and with a suffix of ".yml".

E.g. if your program is called _program_ and you want to run in a production environment, then you would start it using _./program -env=prod_ and your configuration file should be named _program-prod.yml_ and located under _/usr/local/conf_

In your program, if you need to access any configuration element, first save the configuration in a variable, using _conf.GetConf()_.

You can get access to any element in the configuration using its path. You will use _FindNode(conf.Conf,"element.subelement")_ for that purpose. _FindNode_ returns an _interface{}_ and an _error_.

If you want to save yourself the hard work dealing with _interface{}_ elements, you can use 3 helpers:
- _GetString(path string) (string,error)_ will return the _string_ located at _path_
- _GetList(interface{}) ([]yaml.Node,error)_ will return you an array of Yaml Nodes (not necessarily better to deal with)
- _GetStringsMap(path string) (map[string]string,error)_ will return you a _map_ of _string_ elements corresponding to the _path_ informed.


