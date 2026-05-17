import sys

def merge_sort(arr):
    # Base case: arrays with 0 or 1 element are already sorted.
    if len(arr) <= 1:
        return arr
    
    # Split the array into two halves.
    mid = len(arr) // 2
    left = merge_sort(arr[:mid])
    right = merge_sort(arr[mid:])
    
    # Merge the sorted halves.
    return merge(left, right)

def merge(left, right):
    result = []
    i = j = 0
    
    # Compare elements from both lists and merge them in sorted order.
    while i < len(left) and j < len(right):
        if left[i] < right[j]:
            result.append(left[i])
            i += 1
        else:
            result.append(right[j])
            j += 1
            
    # Append any remaining elements.
    result.extend(left[i:])
    result.extend(right[j:])
    return result

def run_tests():
    test_cases = [
        ([], []),
        ([1], [1]),
        ([3, 1, 4, 1, 5, 9, 2, 6], [1, 1, 2, 3, 4, 5, 6, 9]),
        ([9, 8, 7, 6, 5], [5, 6, 7, 8, 9]),
        ([-5, 0, 5, -10], [-10, -5, 0, 5]),
        ([2, 2, 2, 2], [2, 2, 2, 2])
    ]
    
    print("Running Merge Sort Tests...")
    print("-" * 40)
    all_passed = True
    for i, (input_arr, expected) in enumerate(test_cases):
        result = merge_sort(input_arr)
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
    # Check if the correct flag is passed
    if len(sys.argv) < 2 or sys.argv[1] != "--sort":
        # Intentionally cause an error if the flag is missing or incorrect
        raise RuntimeError("Missing or invalid arguments.\n\n"
                           "USAGE ERROR: You must run this script with the '--sort' flag.\n"
                           "Example: python3 merge_sort.py --sort")
    else:
        run_tests()