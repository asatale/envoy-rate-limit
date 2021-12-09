import logging
import sys

stdout_handler = logging.StreamHandler(sys.stdout)
logging.basicConfig(
    level=logging.INFO,
    format='[%(asctime)s] [%(name)s] [%(filename)s:%(lineno)d] %(levelname)s: %(message)s',
    handlers=[stdout_handler])

logger = logging.getLogger("grpc-server")
