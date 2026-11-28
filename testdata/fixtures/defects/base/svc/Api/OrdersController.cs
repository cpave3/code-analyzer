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
        catch (OverflowException e)
        {
            throw new InvalidOperationException("total overflow", e);
        }
    }

    public string Describe(decimal[] lines)
    {
        using var scope = Logging.Begin("describe");
        return Get(lines).ToString("0.00");
    }
}
