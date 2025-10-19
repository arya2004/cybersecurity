# To run this script, open a terminal in this folder and run:
# python3 caesar_crack.py

import collections

# Standard letter frequencies in the English language (in percentages)
# Source: Wikipedia
ENGLISH_FREQUENCIES = {
    'E': 12.70, 'T': 9.06, 'A': 8.17, 'O': 7.51, 'I': 6.97, 'N': 6.75,
    'S': 6.33, 'H': 6.09, 'R': 5.99, 'D': 4.25, 'L': 4.03, 'C': 2.78,
    'U': 2.76, 'M': 2.41, 'W': 2.36, 'F': 2.23, 'G': 2.02, 'Y': 1.97,
    'P': 1.93, 'B': 1.29, 'V': 0.98, 'K': 0.77, 'J': 0.15, 'X': 0.15,
    'Q': 0.10, 'Z': 0.07
}

def caesar_decrypt(ciphertext, shift):
    """
    Decrypts a Caesar ciphered text with a given shift.

    Args:
        ciphertext (str): The encrypted message.
        shift (int): The key (number of positions to shift back).

    Returns:
        str: The decrypted plaintext.
    """
    decrypted_text = ""
    for char in ciphertext:
        if 'a' <= char <= 'z':
            # Handle lowercase letters
            shifted_char = chr(((ord(char) - ord('a') - shift + 26) % 26) + ord('a'))
            decrypted_text += shifted_char
        elif 'A' <= char <= 'Z':
            # Handle uppercase letters
            shifted_char = chr(((ord(char) - ord('A') - shift + 26) % 26) + ord('A'))
            decrypted_text += shifted_char
        else:
            # Keep non-alphabetic characters (spaces, punctuation) as they are
            decrypted_text += char
    return decrypted_text

def score_text(text):
    """
    Scores a text based on how closely its letter frequencies match standard English.
    A lower score indicates a better match.

    Args:
        text (str): The text to score.

    Returns:
        float: The calculated frequency score.
    """
    # 1. Calculate frequencies of letters in the input text
    text_upper = text.upper()
    total_letters = sum(1 for char in text_upper if 'A' <= char <= 'Z')
    if total_letters == 0:
        return float('inf') # Avoid division by zero for empty/non-alphabetic strings

    # Count occurrences of each letter
    letter_counts = collections.Counter(char for char in text_upper if 'A' <= char <= 'Z')

    text_frequencies = {
        char: (count / total_letters) * 100
        for char, count in letter_counts.items()
    }

    # 2. Compare frequencies and calculate score
    # The score is the sum of the absolute differences from English frequencies.
    score = 0.0
    for char_code in range(ord('A'), ord('Z') + 1):
        char = chr(char_code)
        # Get the frequency from our text, or 0.0 if the letter isn't present
        text_freq = text_frequencies.get(char, 0.0)
        # Get the standard English frequency
        english_freq = ENGLISH_FREQUENCIES.get(char, 0.0)
        score += abs(text_freq - english_freq)

    return score

def brute_force_crack(ciphertext):
    """
    Tries all 26 possible keys to decrypt the ciphertext and prints each result.
    """
    print("--- Caesar Cipher Brute Force Attack ---")
    for key in range(26):
        decrypted_text = caesar_decrypt(ciphertext, key)
        print(f"Key #{key:02d}: {decrypted_text}")

def frequency_analysis_crack(ciphertext):
    """
    Uses frequency analysis to find the most likely key and decryption.
    """
    print("\n--- Frequency Analysis Attack ---")
    best_guess = {
        'key': -1,
        'score': float('inf'),
        'text': ''
    }

    for key in range(26):
        decrypted_text = caesar_decrypt(ciphertext, key)
        current_score = score_text(decrypted_text)

        if current_score < best_guess['score']:
            best_guess['score'] = current_score
            best_guess['key'] = key
            best_guess['text'] = decrypted_text

    print(f"Best guess found with Key #{best_guess['key']}:")
    print(f"Decrypted Text: {best_guess['text']}")
    print(f"(Frequency Score: {best_guess['score']:.2f})")


def main():
    """Main function to run the cracking script."""
    # A sample ciphertext encrypted with a key of 15 (shift of 'P')
    ciphertext = "Ymj vznhp gwtbs ktc ozrux tajw ymj qfed itl."

    print(f"Ciphertext to crack: \"{ciphertext}\"\n")

    # Method 1: Brute Force
    brute_force_crack(ciphertext)

    # Method 2: Frequency Analysis
    frequency_analysis_crack(ciphertext)

if __name__ == "__main__":
    main()