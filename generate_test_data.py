import random
import sys

TOTAL_IPS = 30000
DUPLICATE_FACTOR = 0.01

def generate_ip():
    return ".".join(str(random.randint(0, 255)) for _ in range(4))

def generate_ip_list(total_ips, duplicate_factor):
    unique_count = int(total_ips * (1 - duplicate_factor))
    duplicate_count = total_ips - unique_count

    unique_ips = set()
    while len(unique_ips) < unique_count:
        unique_ips.add(generate_ip())

    unique_ips = list(unique_ips)
    duplicates = [random.choice(unique_ips) for _ in range(duplicate_count)]

    ip_list = unique_ips + duplicates
    random.shuffle(ip_list)
    return ip_list

def save_to_file(ip_list, filename):
    with open(filename, 'w') as file:
        for ip in ip_list:
            file.write(ip + '\n')

if __name__ == "__main__":
    if len(sys.argv) > 1:
        total_ips = int(sys.argv[1])
    else:
        total_ips = TOTAL_IPS

    ip_addresses = generate_ip_list(total_ips, DUPLICATE_FACTOR)
    save_to_file(ip_addresses, 'resources/generated_ip_addresses.txt')
