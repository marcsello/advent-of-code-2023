# coded on my phone

def load():
	cmap = []
	with open('input.txt', 'r') as f:
		for i, line in enumerate(f):
			_, lp = line.split(":")
			winnings, haves = lp.split("|")
			winnings = winnings.strip().replace('  ',' ')
			haves = haves.strip().replace('  ',' ')
			winning = set(map(int, winnings.split(' ')))
			have = set(map(int, haves.split(' ')))
			matches = len(winning.intersection(have))
			cmap.append((i, matches))

	return cmap



def main():
	cmap = load()
	nums = [1]*len(cmap)
	for i, score in cmap:
		for j in range(score):
			nums[i+j+1] += nums[i]
	print(sum(nums))

main()
