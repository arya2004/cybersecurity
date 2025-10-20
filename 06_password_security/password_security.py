"""
Password Security Tool - Comprehensive Password Analyzer and Generator
Author: vatsalgupta2004
Description: Analyzes password strength, detects vulnerabilities, and generates secure passwords
"""

import re
import math
import secrets
import string
from typing import Dict, List, Tuple, Set

# Common weak passwords (truncated list for demonstration)
COMMON_PASSWORDS = {
    "password", "123456", "123456789", "12345678", "12345", "1234567", "password1",
    "123123", "1234567890", "qwerty", "abc123", "111111", "1234", "admin", "letmein",
    "welcome", "monkey", "dragon", "master", "sunshine", "princess", "login", "admin123",
    "solo", "1q2w3e4r", "starwars", "qwertyuiop", "654321", "batman", "superman"
}

# Keyboard patterns
KEYBOARD_PATTERNS = [
    "qwerty", "asdf", "zxcv", "qazwsx", "!qaz", "1qaz", "zaq1",
    "123456", "098765", "asdfgh", "qwertyuiop", "zxcvbnm"
]

# Common words to detect
COMMON_WORDS = {"password", "pass", "admin", "user", "root", "login", "welcome"}


class PasswordAnalyzer:
    """Analyzes password strength and vulnerabilities"""
    
    def __init__(self, password: str):
        self.password = password
        self.length = len(password)
        
    def calculate_entropy(self) -> float:
        """
        Calculate Shannon entropy of the password
        Higher entropy = more random and secure
        """
        if not self.password:
            return 0.0
        
        # Count character frequencies
        freq = {}
        for char in self.password:
            freq[char] = freq.get(char, 0) + 1
        
        # Calculate entropy using Shannon's formula
        entropy = 0.0
        for count in freq.values():
            probability = count / self.length
            entropy -= probability * math.log2(probability)
        
        # Normalize by length to get bits of entropy
        return entropy * self.length
    
    def get_character_sets(self) -> Dict[str, bool]:
        """Detect which character sets are used"""
        return {
            'lowercase': bool(re.search(r'[a-z]', self.password)),
            'uppercase': bool(re.search(r'[A-Z]', self.password)),
            'numbers': bool(re.search(r'\d', self.password)),
            'special': bool(re.search(r'[!@#$%^&*()_+\-=\[\]{};:\'",.<>?/\\|`~]', self.password))
        }
    
    def calculate_charset_size(self) -> int:
        """Calculate the size of character set used"""
        charset = self.get_character_sets()
        size = 0
        if charset['lowercase']:
            size += 26
        if charset['uppercase']:
            size += 26
        if charset['numbers']:
            size += 10
        if charset['special']:
            size += 32
        return size
    
    def detect_sequential_patterns(self) -> List[str]:
        """Detect sequential character patterns"""
        patterns = []
        
        # Check for sequential numbers
        if re.search(r'(?:012|123|234|345|456|567|678|789|987|876|765|654|543|432|321|210)', self.password):
            patterns.append("Sequential numbers")
        
        # Check for sequential letters
        lower = self.password.lower()
        if re.search(r'(?:abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz)', lower):
            patterns.append("Sequential letters")
        
        return patterns
    
    def detect_repeated_characters(self) -> List[str]:
        """Detect repeated character patterns"""
        patterns = []
        
        # Check for repeated characters (aaa, 111, etc.)
        if re.search(r'(.)\1{2,}', self.password):
            patterns.append("Repeated characters")
        
        # Check for repeated patterns (abcabc)
        for length in range(2, min(6, self.length // 2 + 1)):
            pattern = self.password[:length]
            if self.password.count(pattern) > 1:
                patterns.append(f"Repeated pattern: '{pattern}'")
                break
        
        return patterns
    
    def detect_keyboard_patterns(self) -> List[str]:
        """Detect keyboard walking patterns"""
        patterns = []
        lower = self.password.lower()
        
        for pattern in KEYBOARD_PATTERNS:
            if pattern in lower:
                patterns.append(f"Keyboard pattern: '{pattern}'")
        
        return patterns
    
    def detect_common_words(self) -> List[str]:
        """Detect common words in password"""
        found_words = []
        lower = self.password.lower()
        
        for word in COMMON_WORDS:
            if word in lower:
                found_words.append(f"Common word: '{word}'")
        
        return found_words
    
    def check_common_password(self) -> bool:
        """Check if password is in common password list"""
        return self.password.lower() in COMMON_PASSWORDS
    
    def calculate_strength_score(self) -> Tuple[int, str]:
        """
        Calculate overall password strength score (0-100)
        Returns: (score, rating)
        """
        score = 0
        
        # Length scoring (0-30 points)
        if self.length >= 16:
            score += 30
        elif self.length >= 12:
            score += 25
        elif self.length >= 8:
            score += 15
        elif self.length >= 6:
            score += 5
        
        # Character variety (0-25 points)
        charset = self.get_character_sets()
        score += sum(charset.values()) * 6.25
        
        # Entropy bonus (0-25 points)
        entropy = self.calculate_entropy()
        if entropy >= 80:
            score += 25
        elif entropy >= 60:
            score += 20
        elif entropy >= 40:
            score += 10
        elif entropy >= 20:
            score += 5
        
        # Deduct points for vulnerabilities (0-20 points)
        vulnerabilities = 0
        if self.check_common_password():
            vulnerabilities += 10
        vulnerabilities += len(self.detect_sequential_patterns()) * 2
        vulnerabilities += len(self.detect_repeated_characters()) * 2
        vulnerabilities += len(self.detect_keyboard_patterns()) * 3
        vulnerabilities += len(self.detect_common_words()) * 2
        
        score = max(0, score - min(vulnerabilities, 20))
        
        # Determine rating
        if score >= 80:
            rating = "VERY STRONG"
        elif score >= 60:
            rating = "STRONG"
        elif score >= 40:
            rating = "MODERATE"
        elif score >= 20:
            rating = "WEAK"
        else:
            rating = "VERY WEAK"
        
        return min(100, score), rating
    
    def estimate_crack_time(self) -> str:
        """Estimate time to crack via brute force"""
        charset_size = self.calculate_charset_size()
        if charset_size == 0:
            return "Instant"
        
        # Assume 10 billion attempts per second (modern GPU)
        attempts_per_second = 10_000_000_000
        total_combinations = charset_size ** self.length
        seconds = total_combinations / attempts_per_second
        
        if seconds < 1:
            return "Instant"
        elif seconds < 60:
            return f"{seconds:.1f} seconds"
        elif seconds < 3600:
            return f"{seconds/60:.1f} minutes"
        elif seconds < 86400:
            return f"{seconds/3600:.1f} hours"
        elif seconds < 31536000:
            return f"{seconds/86400:.1f} days"
        elif seconds < 31536000 * 100:
            return f"{seconds/31536000:.1f} years"
        elif seconds < 31536000 * 1000:
            return f"{seconds/31536000:.0f} years"
        elif seconds < 31536000 * 1_000_000:
            return f"{seconds/(31536000*1000):.1f} thousand years"
        elif seconds < 31536000 * 1_000_000_000:
            return f"{seconds/(31536000*1_000_000):.1f} million years"
        else:
            return f"{seconds/(31536000*1_000_000_000):.1f} billion years"
    
    def get_full_analysis(self) -> Dict:
        """Get comprehensive password analysis"""
        score, rating = self.calculate_strength_score()
        charset = self.get_character_sets()
        
        return {
            'password': self.password,
            'length': self.length,
            'score': score,
            'rating': rating,
            'entropy': self.calculate_entropy(),
            'charset': charset,
            'charset_size': self.calculate_charset_size(),
            'common_password': self.check_common_password(),
            'sequential_patterns': self.detect_sequential_patterns(),
            'repeated_characters': self.detect_repeated_characters(),
            'keyboard_patterns': self.detect_keyboard_patterns(),
            'common_words': self.detect_common_words(),
            'crack_time': self.estimate_crack_time()
        }


class PasswordGenerator:
    """Generates cryptographically secure passwords"""
    
    @staticmethod
    def generate_password(
        length: int = 16,
        use_uppercase: bool = True,
        use_lowercase: bool = True,
        use_numbers: bool = True,
        use_special: bool = True,
        exclude_ambiguous: bool = False
    ) -> str:
        """Generate a random secure password"""
        
        if length < 4:
            raise ValueError("Password length must be at least 4 characters")
        
        # Build character set
        chars = ""
        if use_lowercase:
            chars += string.ascii_lowercase
        if use_uppercase:
            chars += string.ascii_uppercase
        if use_numbers:
            chars += string.digits
        if use_special:
            chars += "!@#$%^&*()_+-=[]{}|;:,.<>?"
        
        if not chars:
            raise ValueError("At least one character set must be selected")
        
        # Exclude ambiguous characters if requested
        if exclude_ambiguous:
            ambiguous = "O0l1I|"
            chars = ''.join(c for c in chars if c not in ambiguous)
        
        # Generate password ensuring at least one char from each selected set
        password = []
        
        # Add at least one character from each selected set
        if use_lowercase:
            password.append(secrets.choice(string.ascii_lowercase))
        if use_uppercase:
            password.append(secrets.choice(string.ascii_uppercase))
        if use_numbers:
            password.append(secrets.choice(string.digits))
        if use_special:
            password.append(secrets.choice("!@#$%^&*()_+-=[]{}|;:,.<>?"))
        
        # Fill the rest randomly
        for _ in range(length - len(password)):
            password.append(secrets.choice(chars))
        
        # Shuffle to avoid predictable patterns
        secrets.SystemRandom().shuffle(password)
        
        return ''.join(password)
    
    @staticmethod
    def generate_passphrase(num_words: int = 4, separator: str = "-") -> str:
        """Generate a memorable passphrase"""
        # Simple word list for demonstration
        words = [
            "alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "hotel",
            "india", "juliet", "kilo", "lima", "mike", "november", "oscar", "papa",
            "quebec", "romeo", "sierra", "tango", "uniform", "victor", "whiskey", "xray",
            "yankee", "zulu", "cipher", "encrypt", "secure", "shield", "guard", "protect",
            "fortress", "vault", "lock", "key", "token", "badge", "sign", "mark"
        ]
        
        selected_words = [secrets.choice(words) for _ in range(num_words)]
        return separator.join(selected_words)


def print_analysis(analysis: Dict):
    """Pretty print password analysis"""
    print("\n" + "="*60)
    print(f"🔍 Password Analysis")
    print("="*60)
    
    print(f"\n📊 Strength Score: {analysis['score']}/100 - {analysis['rating']}")
    print(f"🔐 Entropy: {analysis['entropy']:.2f} bits")
    print(f"📏 Length: {analysis['length']} characters")
    print(f"🎲 Character Set Size: {analysis['charset_size']}")
    print(f"⏱️  Estimated Crack Time: {analysis['crack_time']}")
    
    print(f"\n✓ Character Types:")
    for charset, present in analysis['charset'].items():
        symbol = "✅" if present else "❌"
        print(f"  {symbol} {charset.capitalize()}")
    
    if analysis['common_password']:
        print(f"\n⚠️  WARNING: This is a commonly used password!")
    
    vulnerabilities = []
    vulnerabilities.extend(analysis['sequential_patterns'])
    vulnerabilities.extend(analysis['repeated_characters'])
    vulnerabilities.extend(analysis['keyboard_patterns'])
    vulnerabilities.extend(analysis['common_words'])
    
    if vulnerabilities:
        print(f"\n⚠️  Vulnerabilities Found:")
        for vuln in vulnerabilities:
            print(f"  ⚠️  {vuln}")
    else:
        print(f"\n✅ No common vulnerabilities detected!")
    
    print("="*60)


def main():
    """Main interactive CLI"""
    print("="*60)
    print("🔐 Password Security Tool")
    print("="*60)
    print("\nOptions:")
    print("1. Analyze password strength")
    print("2. Generate secure password")
    print("3. Generate passphrase")
    print("4. Batch generate passwords")
    print("5. Exit")
    
    while True:
        choice = input("\nSelect option (1-5): ").strip()
        
        if choice == "1":
            password = input("\nEnter password to analyze: ")
            analyzer = PasswordAnalyzer(password)
            analysis = analyzer.get_full_analysis()
            print_analysis(analysis)
        
        elif choice == "2":
            try:
                length = int(input("Password length (8-128): ") or "16")
                password = PasswordGenerator.generate_password(length=length)
                print(f"\n🎲 Generated Password: {password}")
                
                # Analyze generated password
                analyzer = PasswordAnalyzer(password)
                score, rating = analyzer.calculate_strength_score()
                print(f"Strength: {score}/100 - {rating}")
            except ValueError as e:
                print(f"❌ Error: {e}")
        
        elif choice == "3":
            try:
                num_words = int(input("Number of words (3-8): ") or "4")
                passphrase = PasswordGenerator.generate_passphrase(num_words=num_words)
                print(f"\n🎲 Generated Passphrase: {passphrase}")
                
                # Analyze generated passphrase
                analyzer = PasswordAnalyzer(passphrase)
                score, rating = analyzer.calculate_strength_score()
                print(f"Strength: {score}/100 - {rating}")
            except ValueError as e:
                print(f"❌ Error: {e}")
        
        elif choice == "4":
            try:
                count = int(input("How many passwords to generate (1-20): ") or "5")
                length = int(input("Password length (8-128): ") or "16")
                print(f"\n🎲 Generated {count} Passwords:")
                for i in range(min(count, 20)):
                    password = PasswordGenerator.generate_password(length=length)
                    print(f"  {i+1}. {password}")
            except ValueError as e:
                print(f"❌ Error: {e}")
        
        elif choice == "5":
            print("\n👋 Goodbye!")
            break
        
        else:
            print("❌ Invalid option. Please select 1-5.")


if __name__ == "__main__":
    main()
