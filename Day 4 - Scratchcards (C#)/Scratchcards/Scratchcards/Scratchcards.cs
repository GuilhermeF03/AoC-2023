using System.Text.RegularExpressions;

namespace Scratchcards
{
    internal class Card
    {
        internal readonly int Idx; 
        internal readonly List<int> WinningNumbers;
        internal readonly List<int> MyNumbers;
        
        public Card(int idx, List<int> winningNumbers, List<int> myNumbers)
        {
            Idx = idx;
            WinningNumbers = winningNumbers;
            MyNumbers = myNumbers;
        }
    }

    internal class CardPointer
    {
        internal readonly int Id;
        public CardPointer(int id) { Id = id; }
    }

    internal abstract partial class Scratchcards
    {
        // Regex patterns
        [GeneratedRegex("\\d+")]
        private static partial Regex NumberRegex();
        [GeneratedRegex( @"(Card)(\s+\d+):")]
        private static partial Regex CardHeaderRegex();
        
        private const string CardHeaderPatternRegex = @"(Card)(\s+\d+):";
        private const string WinningNumberPatternRegex = @"(\d+\s+)+\|";
        private const string MyNumbersPatternRegex = @"\|(\s+\d+)+";
        
        private static readonly LinkedList<CardPointer> CardPointers = new();
        private static readonly List<Card> Cards = new();
        private const string FilePath = "../../../input.txt";

        private static void Main()
        {
            var lines = File.ReadLines(FilePath).ToList();
            Console.WriteLine("Part1: {0}", Part1(lines));
            Console.WriteLine("Part2: {0}", Part2(lines));
        }
        
        private static IEnumerable<int> GetParsedNumbers(string input, string parseRegex)
        {
            var parsedLine = Regex.Matches(input, parseRegex).First().Value;
            return NumberRegex().Matches(parsedLine).Select(match => int.Parse(match.Value));
        }
        
        private static int Part1(IEnumerable<string> lines)
        {
            return (
                from line in lines
                    let header = CardHeaderRegex().Match(line).Value
                    let winningNumbers = GetParsedNumbers(line, WinningNumberPatternRegex)
                    let myNumbers = GetParsedNumbers(line, MyNumbersPatternRegex)
                select winningNumbers.Intersect(myNumbers).ToList() into matchedNumbers 
                select (int)Math.Pow(2, matchedNumbers.Count - 1)
                ).Sum();
        }

        private static int Part2(IEnumerable<string> lines)
        {
            // fill list with initial values
            foreach (var line in lines)
            {
                var card = new Card(
                    idx: GetParsedNumbers(line, CardHeaderPatternRegex).First(),
                    winningNumbers: GetParsedNumbers(line, WinningNumberPatternRegex).ToList(),
                    myNumbers: GetParsedNumbers(line, MyNumbersPatternRegex).ToList()
                );
                CardPointers?.AddLast(new CardPointer(card.Idx - 1));
                Cards.Add(card);
            }
            
            // iterate over pointers
            if (CardPointers == null) return 0;
            {
                var currNode = CardPointers.First;
                while (currNode?.Next != null)
                {
                    var cardPointer = currNode.Value;
                    // check how many intersections are there
                    var card = Cards[cardPointer.Id];
                    var matchedNumbersSize = card.WinningNumbers.Intersect(card.MyNumbers).Count();
                    // insert copies
                    for (var i = 0; i < matchedNumbersSize; i++)
                        if (currNode != null && cardPointer.Id + i + 1 < Cards.Count)
                            CardPointers.AddAfter(currNode, new CardPointer(cardPointer.Id + i + 1));
                    currNode = currNode?.Next;
                }
            }
            return CardPointers.Count;
        }

    }
}
