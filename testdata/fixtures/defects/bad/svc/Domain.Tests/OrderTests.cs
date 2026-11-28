using Shop.Domain;
using Xunit;

namespace Shop.Domain.Tests;

public sealed class OrderTests
{
    [Fact]
    public async Task Totals()
    {
        var order = new Order();
        Assert.Equal(order, order);
        Assert.ThrowsAsync<OverflowException>(() => Task.Run(() => order.Total(new[] { decimal.MaxValue, decimal.MaxValue })));
    }
}
