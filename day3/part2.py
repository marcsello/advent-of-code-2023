import string

TOTALLY_NOT_A_SYMBOL = string.digits + '.'


def magic_loader() -> tuple[list[list], list]:
    table = []
    numbers_lut = []
    meme_id = 0  # This unique id will have a significant role later on
    with open("input.txt", "r") as f:
        row = []
        current_int = None
        for line in f:
            row = []
            line = line.strip()  # remove the \n
            current_int = None
            for i, c in enumerate(line):
                # iterate over chars
                if c in string.digits:
                    if current_int is None:
                        # found the first character of an int
                        # search for the end of the number
                        end_i = i
                        while True:
                            end_i += 1
                            if end_i == len(line):
                                # the char after
                                break
                            if line[end_i] not in string.digits:
                                break

                        current_int = int(line[i:end_i])
                        assert len(numbers_lut) == meme_id, f"oopsie woopsie, the array is broken at index {meme_id}"
                        numbers_lut.append(current_int)

                else:
                    # not a number
                    if current_int:
                        meme_id += 1
                        current_int = None

                row.append((c, current_int, meme_id if current_int else None))

            assert len(row) == len(line), "line is not good length"
            table.append(row)
            if current_int:  # work around numbers ending at the line end
                meme_id += 1

    return table, numbers_lut


def get_surroundings(rowi, coli, table):  # generator
    tests = [  # clockwise
        (rowi, coli + 1),
        (rowi + 1, coli + 1),
        (rowi+1, coli),
        (rowi + 1, coli - 1),
        (rowi, coli - 1),
        (rowi - 1, coli - 1),
        (rowi-1, coli),
        (rowi-1, coli+1),
    ]
    for trowi, tcoli in tests:
        try:
            yield table[trowi][tcoli]
        except IndexError:
            pass


def main():
    # load data into a highly advanced data structure
    table, numbers_lut = magic_loader()

    # collect all numbers surrounding the symbols
    summary = 0
    for rowi, row in enumerate(table):
        for coli, col in enumerate(row):
            if col[0] == '*':
                # is a gear
                neighbors = set()
                for _, _, meme_id in get_surroundings(rowi, coli, table):
                    if meme_id is not None:
                        neighbors.add(meme_id)

                if len(neighbors) == 2:
                    nl = list(neighbors)
                    summary += numbers_lut[nl[0]] * numbers_lut[nl[1]]

    # summarize
    print(summary)


if __name__ == '__main__':
    main()
