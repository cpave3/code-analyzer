namespace Shop.Domain;

public sealed class Order
{
    public decimal Total(decimal[] lines)
    {
        decimal sum = 0;
        foreach (var line in lines)
        {
            sum += Round(line);
        }
        return sum;
    }

    private static decimal Round(decimal value) => decimal.Round(value, 2);
}
