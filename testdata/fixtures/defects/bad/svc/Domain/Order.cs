using Shop.Api;

namespace Shop.Domain;

public sealed class Order
{
    public decimal Total(decimal[] lines)
    {
        decimal sum = 0;
        var unusedCount = lines.Length;
        foreach (var line in lines)
        {
            sum += Round(line);
        }
        return sum;
        Audit(sum);
    }

    private static decimal Round(decimal value) => decimal.Round(value, 2);

    private static void LegacyRound(decimal value) { }

    private static void Audit(decimal value) => OrdersController.Notify(value);
}
