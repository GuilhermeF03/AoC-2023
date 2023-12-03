import java.lang.IllegalStateException
import java.nio.file.Path
import kotlin.io.path.readLines

val redPattern = Regex("(\\d+) (red)")
val greenPattern = Regex("(\\d+) (green)")
val bluePattern = Regex("(\\d+) (blue)")
val gameIdPattern = Regex("(Game) (\\d+)(:)")



val lines = Path.of("./input.txt").readLines()
fun main (){
    println(part1())
    println(part2())
}

data class GameResult(val reds :Int, val greens : Int, val blues :Int)
data class GameInfo(val id: Int, val results : List<GameResult>)

fun part1() : Int =
    lines.asSequence()
        .map { it.split(";") }
        .map{parseGameInfo(it)}
        .filter {
            it.results.all { game ->
                game.reds <= 12 && game.greens <= 13 && game.blues <= 14
            }
        }
        .sumOf {it.id}

fun parseGameInfo(gameLine : List<String>): GameInfo {
    val id : Int = gameIdPattern.find(gameLine[0])?.value?.parseInt()
            ?: throw IllegalStateException("No Id was found")

    val reds = gameLine.map { redPattern.findAll(it).mergeInts() }
    val greens = gameLine.map { greenPattern.findAll(it).mergeInts() }
    val blues = gameLine.map { bluePattern.findAll(it).mergeInts() }

    return GameInfo(id, List(gameLine.size) { idx -> GameResult(reds[idx], greens[idx], blues[idx]) })
}


fun Sequence<MatchResult>.mergeInts() : Int =
        this.map { it.groupValues[0] }.toList().sumOf { it.parseInt() }
fun String.parseInt(): Int = this.filter { it.isDigit() }.toInt()

fun part2() : Int =
    lines.asSequence()
        .map { it.split(";") }
        .map{
            val gameInfo = parseGameInfo(it)
            findMinimumSet(gameInfo)
        }
        .sum()


fun findMinimumSet(gameInfo : GameInfo) : Int {
    val minimumRed = gameInfo.results.map { it.reds }.maxOf { it }
    val minimumGreen = gameInfo.results.map { it.greens }.maxOf { it }
    val minimumBlue = gameInfo.results.map { it.blues }.maxOf { it }
    return minimumRed * minimumGreen * minimumBlue
}