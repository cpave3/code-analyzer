using Shop.Domain;

namespace Shop.Api;

public sealed class OrdersController
{
    private readonly Order _order = new();

    public decimal Get(decimal[] lines)
    {
        try
        {
            return _order.Total(lines);
        }
        catch (OverflowException)
        {
        }
        return 0;
    }

    public string Describe(decimal[] lines)
    {
        try
        {
            return Get(lines).ToString("0.00");
        }
        catch (FormatException)
        {
            return null;
        }
    }
}
