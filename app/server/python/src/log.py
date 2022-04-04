import logging
import sys
from config import cfg

log_levels =  {
    "DEBUG": logging.DEBUG,
    "INFO": logging.INFO,
    "WARN": logging.WARN,
    "ERROR": logging.ERROR,
    "CRITICAL": logging.CRITICAL,
    "FATAL": logging.FATAL,
    }

stdout_handler = logging.StreamHandler(sys.stdout)
logging.basicConfig(
    level=log_levels[cfg.log.upper()],
    format='[%(asctime)s] [%(name)s] [%(filename)s:%(lineno)d] %(levelname)s: %(message)s',
    handlers=[stdout_handler])

logger = logging.getLogger("grpc-server")
