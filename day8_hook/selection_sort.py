import sys


def selection_sort(arr: list) -> list:
    result = list(arr)
    n = len(result)
    for i in range(n):
        min_idx = i
        for j in range(i + 1, n):
            if result[j] < result[min_idx]:
                min_idx = j
        result[i], result[min_idx] = result[min_idx], result[i]
    return result


def run_tests() -> None:
    test_cases = [
        ([], []),
        ([1], [1]),
        ([3, 1, 4, 1, 5, 9, 2, 6], [1, 1, 2, 3, 4, 5, 6, 9]),
        ([9, 8, 7, 6, 5], [5, 6, 7, 8, 9]),
        ([-5, 0, 5, -10], [-10, -5, 0, 5]),
        ([2, 2, 2, 2], [2, 2, 2, 2]),
    ]

    print("Running Selection Sort Tests...")
    print("-" * 40)
    all_passed = True
    for i, (input_arr, expected) in enumerate(test_cases):
        result = selection_sort(input_arr)
        if result == expected:
            print(f"Test {i+1} PASSED: {input_arr} -> {result}")
        else:
            print(f"Test {i+1} FAILED: Expected {expected}, got {result}")
            all_passed = False

    print("-" * 40)
    if all_passed:
        print("SUCCESS: All tests passed!")
    else:
        print("FAILURE: Some tests failed.")


if __name__ == "__main__":
    if len(sys.argv) < 2 or sys.argv[1] != "--sort":
        raise RuntimeError(
            "Missing or invalid arguments.\n\n"
            "USAGE ERROR: You must run this script with the '--sort' flag.\n"
            "Example: python3 selection_sort.py --sort"
        )
    run_tests()
