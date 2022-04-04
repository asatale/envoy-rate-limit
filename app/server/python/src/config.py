import argparse


parser = argparse.ArgumentParser(description='Python GRPC server')

parser.add_argument('-log', type=str, default='INFO',
                    help='Logging level for service. Options: [DEBUG, INFO, WARN, CRITICAL, FATAL, ERROR  ]')
parser.add_argument('-addr', type=str, default='0.0.0.0:50051',
                    help='Server bind address. Default: 0.0.0.0:50051')
parser.add_argument('-delay', type=int, default=10,
                    help='Response delay in millisecond. Default: 10ms')
parser.add_argument('-dprob', type=int, default=20,
                    help='Delay Probability')
parser.add_argument('-cancel', dest='cancel', action='store_const',
                    const=True, default=False,
                    help='Cancel RPC with cancel-probability')
parser.add_argument('-cprob', type=int, default=20,
                    help='Cancel Probability')
parser.add_argument('-metric_addr', type=str, default='0.0.0.0:8000',
                    help='Prometheus scrapping endpoint for service')

cfg = parser.parse_args()
