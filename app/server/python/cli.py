import argparse


parser = argparse.ArgumentParser(description='Python GRPC server')

parser.add_argument('-addr', type=str, default='0.0.0.0:50051',
                    help='Server bind address. Default: 0.0.0.0:50051')
parser.add_argument('-rsp_delay', type=int, default=10,
                    help='Response delay in millisecond. Default: 10ms')
parser.add_argument('-variance', type=int, default=2,
                    help='Response time randomized variance in millisecond. Default: 2ms')

args = parser.parse_args()
