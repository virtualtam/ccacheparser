ccacheparser
============

A simple utility to convert `ccache`_ statistics to a JSON object.

====== ======
Branch Status
====== ======
master .. image:: https://travis-ci.org/virtualtam/ccacheparser.svg?branch=master
          :target: https://travis-ci.org/virtualtam/ccacheparser
          :alt: Travis build status
====== ======

Installation
------------

::

   $ go get -u github.com/virtualtam/ccacheparser


Usage
-----

::

   $ ccache -s | ccacheparser | jq

   {
     "cache_directory": "/home/virtualtam/.ccache",
     "primary_config": "/home/virtualtam/.ccache/ccache.conf",
     "secondary_config_readonly": "/etc/ccache.conf",
     "stats_time": "2018-09-24T21:19:07.997866938+02:00",
     "stats_zero_time": "2018-09-23T01:18:52+02:00",
     "cache_hit_direct": 124,
     "cache_hit_preprocessed": 8,
     "cache_miss": 297,
     "cache_hit_rate": 30.77,
     "called_for_link": 39,
     "called_for_preprocessing": 263,
     "unsupported_code_directive": 5,
     "no_input_file": 83,
     "cleanups_performed": 0,
     "files_in_cache": 926,
     "cache_size": "17.5 MB",
     "cache_size_bytes": 17500000,
     "max_cache_size": "15.0 GB",
     "max_cache_size_bytes": 15000000000
   }


.. _ccache: https://github.com/ccache/ccache
