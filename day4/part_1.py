# coded on my phone

def main():
	sum_ = 0
	with open('input.txt', 'r') as f:
		for line in f:
			_, lp = line.split(":")
			winnings, haves = lp.split("|")
			winnings = winnings.strip().replace('  ',' ')
			haves = haves.strip().replace('  ',' ')
			winning = set(map(int, winnings.split(' ')))
			have = set(map(int, haves.split(' ')))
			matches = len(winning.intersection(have))
			score = 0
			if matches > 0:
				score = 2 ** (matches-1)
			sum_ += score

	print(sum_)


main()
