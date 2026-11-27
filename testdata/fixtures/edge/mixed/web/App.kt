fun greet(name: String): String {
    val n = name.trim()
    return if (n.isNotEmpty()) "hello $n" else "hello"
}
