[![Edwig logo](https://github.com/af83/edwig/wiki/images/edwig_logo.png)](https://enroute.mobi/produits/edwig/)

[![Maintainability](https://api.codeclimate.com/v1/badges/bdf4ce25da411be47722/maintainability)](https://codeclimate.com/github/af83/edwig/maintainability)

## :warning: Project changes

Edwig changed name to Ara, and is now located here: https://bitbucket.org/enroute-mobi/ara

## An innovative and modular solution

* Modular architecture organized in Collection, Model & Broadcast
* Multi-protocol connectors: SIRI, SIRI Lite (GTFS-RT soon)
* Real time Visualization / management of data by API
* Loading theoretical offer and / or network structure into a database
* Real time and parameterizable logging
* Managing multiple independent referentials in the same server
* Real time administration: exchange partners, referentials

## SIRI connectors

In collection and broadcast (both subscription and request)

* StopMonitoring
* EstimatedTimeTable
* Situational Management

In broadcast only:

* StopPointDiscovery
* LineDiscovery

## Versatile and multilingual

* Transcodification of data with use and correspondence between different types of identifiers on the same objects
* Management and configuration of identifier formats to adapt in real time the identifiers used with an exchange partner
* Modular import supply by a new product "Referentials"

## Real-time logging

Outsource, process and store in real time all exchanges managed by Edwig :

* Send real time exchange data to LogStash processing
* High performance storage for consultation and statistics via ElasticSearch
* Visualization of historical data via Kibana

## More Information

Some technical articles are available [on the wiki](../../wiki) too.

Related projects :

* [Documentation](https://github.com/af83/edwig-docs)
* [Administration interface](https://github.com/af83/edwig-admin)
* [Ruby SDK](https://github.com/af83/edwig)

## License

This project is licensed under the Apache2 license, a copy of which can be found in the [LICENSE](./LICENSE.md) file.

## Support

Contact [af83 Edwig team](mailto:edwig-dev@af83.com) to know how to contribute to the Edwig project
