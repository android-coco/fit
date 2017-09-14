package fit

/// command line option (all command line options should be support in configuration file):
// --config(*) : setting which the configuration file will be used. find the 'fit.conf' under current directory by default
// --loglevel(*) :  log level.  'none','Info', 'warn','error', 'fatal',
//			  log should be printed on stdout/terminal, file according to the configuration
//

//configuration file option
//DocRoot : where is the web root folder. default is the './'
//LogLevel : log level
//LogFile : 'stdout', file system file
//Port : default is 80
//ReadTimeout : socket read timeout. default is : 10s
//ReadHeaderTimeout: timeout before finish header reading.
//IdleTimeout
//WriteTimeout: socket write timeout. default is : 10s
//MaxHeaderBytes: header size. default is : 1<<20 (1M)
//MaxBytesReader : max incoming data size: default is : 4M
